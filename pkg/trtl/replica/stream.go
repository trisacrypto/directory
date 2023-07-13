package replica

import (
	"io"
	"sync"

	"github.com/rotationalio/honu/replica"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

const streamBufSize = 8

// The gossipStream interface describes both the replica.Replication_GossipClient and
// the replica.Replication_GossipServer without the additional methods defined by grpc.
type gossipStream interface {
	Send(*replica.Sync) error
	Recv() (*replica.Sync, error)
}

// grpc streams have blocking Send and Recv methods that perform networking operations.
// Each method can only be called from one go routine, meaning that one go routine can
// both Send and Recv, or two go routines are required, one for Send and the other for
// Recv. Both the AntiEntropySync and Gossip methods have multiple go routines that need
// to send messages, therefore the streamSender synchronizes these go routines by
// creating a channel they can send messages on, ensuring that only one routine calls
// the stream Send method.
type streamSender struct {
	sync.RWMutex
	wg     *sync.WaitGroup    // A wait group so the caller knows when sending is complete
	ok     bool               // If true, it is ok to send on the channel
	log    *sentry.Logger     // A logger to report send errors on
	msgs   chan *replica.Sync // The channel for go routines to send messages on
	stream gossipStream       // The gossip stream to send messages on
}

// Creates the stream sender and kicks off the sending go routine.
func newStreamSender(wg *sync.WaitGroup, logctx *sentry.Logger, stream gossipStream) *streamSender {
	sender := &streamSender{
		wg:     wg,
		ok:     true,
		log:    logctx,
		msgs:   make(chan *replica.Sync, streamBufSize),
		stream: stream,
	}

	wg.Add(1)
	go sender.Run()
	return sender
}

func (s *streamSender) Run() {
	// Make sure the calling go routine knows when the sender is done.
	defer s.wg.Done()

	// Loop over the messages channel and keep sending messages while it's open.
	for msg := range s.msgs {
		s.log.Trace().Str("status", msg.Status.String()).Msg("sending sync")
		if err := s.stream.Send(msg); err != nil {
			if err != io.EOF {
				s.log.Error().Err(err).Msg("could not send gossip message to peer")
			}

			// Let external go routines know they can no longer send on the channel.
			s.Lock()
			s.ok = false
			s.Unlock()
		}
	}

	// Continue to consume messages until the msgs channel is closed to ensure that
	// go routines trying to send messages don't get blocked if the message channel is
	// filled up and to prevent the race condition between checking ok and sending a
	// message. This won't do anything with the messages.
	//
	// NOTE: ranging over a closed channel immediately returns so even if the msgs
	// channel was closed in the first loop, this second loop will not execute nor
	// will it panic.
	for msg := range s.msgs {
		s.log.Trace().Str("status", msg.Status.String()).Msg("dropped sync message (not sent)")
	}
}

// Send a sync message to the remote by putting it on the msgs channel, returns false if
// the grpc connection has errored to tell the calling go routine to stop trying to send
// messages. Note that just because ok is returned doesn't mean the message was sent
// successfully, it just means that the stream hasn't errored yet.
func (s *streamSender) Send(msg *replica.Sync) bool {
	if s.Ok() {
		s.msgs <- msg
		return true
	}
	return false
}

// Thread-safe method to check if it's ok to send a message. If true is returned it
// indicates that the stream is not currently closed due to an error, it does not mean
// that a message put on the msgs channel will be sent.
func (s *streamSender) Ok() bool {
	s.RLock()
	defer s.RUnlock()
	return s.ok
}

// Close the msgs stream so that external go routines can indicate they're done sending.
// Only one go routine should call the Close method (usually the main go routine).
func (s *streamSender) Close() {
	s.Lock()
	s.ok = false
	s.Unlock()
	close(s.msgs)
}
