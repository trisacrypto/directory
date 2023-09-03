package pagination

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Load a PageCursor from a page_token string.
// This is very similar code to trtl.internal but decoupled so Trtl and GDS aren't
// dependent on each other.
func (pc *PageCursor) Load(token string) (err error) {
	var data []byte
	if data, err = base64.RawURLEncoding.DecodeString(token); err != nil {
		return fmt.Errorf("could not decode page token: %s", err)
	}

	if err = proto.Unmarshal(data, pc); err != nil {
		return fmt.Errorf("could not unmarshal page cursor: %s", err)
	}
	return nil
}

// Dump a PageCursor into a page_token string.
// This is very similar code to trtl.internal but decoupled so Trtl and GDS aren't
// dependent on each other.
func (pc *PageCursor) Dump() (token string, err error) {
	var data []byte
	if data, err = proto.Marshal(pc); err != nil {
		return "", fmt.Errorf("could not marshal page cursor: %s", err)
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
