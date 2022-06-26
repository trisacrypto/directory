package db

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

const (
	NamespaceAnnouncements = "announcements"
)

var (
	lastUpdatedKey = []byte("last_updated")
)

// The Announcements type exposes methods for interacting with Network Announcements
// in the database. Announcements are stored as compact, compressed JSON to minimize the
// network requests and reduce storage requirements. This struct performs all necessary
// serialization on the announcements before storing and retrieving the model object. The
// announcement keys are ksuids - timestamp ordered unique IDs so that it is easy to scan
// the trtl database to find the most recent announcements.
//
// Announcements implements the Collection interface
type Announcements struct {
	db        *DB
	namespace string
}

// Ensure that Announcements implements the Collection interface.
var _ Collection = &Announcements{}

// Announcements constructs the collection type for interactions with network
// announcement objects. This method is intended to be used with chaining, e.g. as
// db.Announcements().Recent(), so to reduce the number of allocations a singleton
// intermediate struct is used. Method calls to the collection are thread-safe.
func (db *DB) Announcements() *Announcements {
	db.muMakeAC.Do(func() {
		db.announcements = &Announcements{
			db:        db,
			namespace: NamespaceAnnouncements,
		}
	})
	return db.announcements
}

// Recent returns the set of results whose post date is after the not before timestamp,
// limited to the maximum number of results. Last updated returns the timestamp that
// any announcement was added or changed.
func (a *Announcements) Recent(ctx context.Context, maxResults int, notBefore time.Time) (out *api.AnnouncementsReply, err error) {
	out = &api.AnnouncementsReply{
		Announcements: make([]*api.Announcement, 0, maxResults),
	}

	// Get the last updated value from the index key in the collection.
	var value []byte
	if value, err = a.db.Get(ctx, lastUpdatedKey, a.namespace); err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("could not fetch last updated: %s", err)
	}

	if len(value) > 0 {
		if err = json.Unmarshal(value, &out.LastUpdated); err != nil {
			return nil, fmt.Errorf("could not unmarshal last updated: %s", err)
		}
	}

	var cursor trtl.Trtl_CursorClient
	if cursor, err = a.db.trtl.Cursor(ctx, &trtl.CursorRequest{Namespace: a.namespace}); err != nil {
		return nil, fmt.Errorf("could not connect to trtl cursor: %s", err)
	}
	defer cursor.CloseSend()

	// Consume the cursor
	// TODO: this filtering method is NOT ideal since it scans in ascending time order
	// rather than reverse time order (trtl doesn't support reverse ordering). This
	// means as the number of announcements in the collection grows, this method will
	// get slower and slower. Worse, we're loading everything in memory so we might
	// overload the pod. We need to either figure out a different key-value storage
	// and ordering mechanism or implement reverse seek in trtl.
	results := make([]*api.Announcement, 0)
	for {
		var pair *trtl.KVPair
		if pair, err = cursor.Recv(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("could not read message from cursor: %s", err)
		}

		if bytes.Equal(pair.Key, lastUpdatedKey) {
			continue
		}

		var post *api.Announcement
		if post, err = a.Decode(pair.Value); err != nil {
			// If we can't decode the announcement log the error and continue
			log.Error().Err(err).Str("key", ParseKSUID(pair.Key)).Msg("could not decode announcement from database")
			continue
		}

		var date time.Time
		if date, err = time.Parse("2006-01-02", post.PostDate); err != nil {
			// If we can't decode the post date log the error and continue
			log.Error().Err(err).Str("key", ParseKSUID(pair.Key)).Msg("could not parse the post date")
		}

		if date.Before(notBefore) {
			// Filter on the not before date
			continue
		}

		results = append(results, post)
	}

	// Reverse the results and apply the max results limit
	for i := len(results) - 1; i >= 0; i-- {
		if len(out.Announcements) >= maxResults {
			break
		}
		out.Announcements = append(out.Announcements, results[i])
	}
	return out, nil
}

// Post an announcement, putting it to the trtl database. This method does no
// verification of duplicate announcements or any content verification except for a
// check that an empty announcement is not being put to the database.
func (a *Announcements) Post(ctx context.Context, in *api.Announcement) (_ string, err error) {
	// Make sure we don't post empty announcements
	if in.Title == "" && in.Body == "" && in.PostDate == "" && in.Author == "" {
		return "", ErrEmptyAnnouncement
	}

	// Encode the Post to bytes for storage
	var value []byte
	if value, err = a.Encode(in); err != nil {
		return "", fmt.Errorf("could not encode announcement: %s", err)
	}

	// Create a ksuid for the post
	key := ksuid.New()

	if err = a.db.Put(ctx, key.Bytes(), value, a.namespace); err != nil {
		return "", fmt.Errorf("could not put announcement to db: %s", err)
	}

	// Set the last updated timestamp on the database
	var updated []byte
	if updated, err = json.Marshal(time.Now()); err != nil {
		return "", fmt.Errorf("could not create last updated timestamp value: %s", err)
	}

	if err = a.db.Put(ctx, lastUpdatedKey, updated, a.namespace); err != nil {
		return "", fmt.Errorf("could not put last updated index: %s", err)
	}

	return key.String(), nil
}

// Namespace implements the collection interface
func (a *Announcements) Namespace() string {
	return a.namespace
}

// Encode an announcement for storage in the database.
// HACK: the api definition is probably not the right place to store a database model,
// should we just create protocol buffer model definitions instead of using this?
func (a *Announcements) Encode(in *api.Announcement) (_ []byte, err error) {
	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)

	if err = json.NewEncoder(w).Encode(in); err != nil {
		return nil, err
	}
	w.Close()

	return buf.Bytes(), nil
}

// Decode an annoucement from storage in the database.
// HACK: the api definition is probably not the right place to store a database model,
// should we just create protocol buffer model definitions instead of using this?
func (a *Announcements) Decode(data []byte) (out *api.Announcement, err error) {
	buf := bytes.NewBuffer(data)

	var r *gzip.Reader
	if r, err = gzip.NewReader(buf); err != nil {
		return nil, err
	}
	defer r.Close()

	if err = json.NewDecoder(r).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// ParseKSUID from bytes into a string, returning the nil ksuid if an error occurs.
func ParseKSUID(key []byte) string {
	kid, err := ksuid.FromBytes(key)
	if err != nil {
		return ksuid.Nil.String()
	}
	return kid.String()
}
