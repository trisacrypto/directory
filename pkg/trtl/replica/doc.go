/*
Package replica implements bilateral anti-entropy replication for the trtl database.

Each trtl node is allowed to respond to accesses by the user independently, e.g. without
coordinating with other trtl nodes. As they apply Put and Delete operations to their
local data store, the databases of each node diverge, causing entropy. Periodically,
each node will initiate an anti-entropy session with another trtl node so that their
local databases are synchronized, thus reducing the entropy caused by divergent accesses
on each node. In the absence of writes, all nodes will eventually become consistent with
each other, meaning that they will all have the same exact copy of the database.

This package implements both the periodic initiation of the anti-entropy session as well
as the remote handler for responding to a synchronization request. Our anti-entropy
implementation is bilateral, meaning that both the initiating replica and the remote
responding replica will be synchronized. Generally speaking there are two types of
unilateral anti-entropy: push and pull (generally referred to as gossip protocols). In
unilateral push anti-entropy, the initiating replica sends its latest versions to the
remote (e.g. the remote is being synchronized); in unilateral pull anti-entropy, the
initiating replica requests later versions from the remote (the initiator is being
synchronized). Both types of unilateral anti-entropy are one- or two-phased; bilateral
anti-entropy requires four (or possibly five, depending on how you count) phases and
relies heavily on gRPC streaming for success. The implementation of anti-entropy in both
sync.go and replica.go is based on these phases, so it's important to understand them
to maintain this library.

Anti-entropy sessions are initiated by a long running go routine, go AntiEntropy, which
acts as the client when connecting to the remote peer. The AntiEntropy routine sleeps
for a jittered interval -- a random amount of time normally distributed by a mean and
standard deviation duration -- to ensure that anti-entropy sessions are not all
happening concurrently. After this interval has passed, the initiator will select a
random peer in the network, and will open up a Gossip stream to it. It then begins two
phases in three go routines to start the anti-entropy session. We refer to these phases
as the initiator phase 1 and phase 2 respectively.

The Gossip stream is handled by the remote node in the Gossip method - the gRPC handler
defined by the Replication Service in protocol buffers. The Gossip handler also runs
three go routines for handling bilateral anti-entropy. It also executes two phases, the
first phase of which starts when the stream is created and the second phase when it
receives a message from the initiator (described below). The first phase is divided into
two parts, the part before the second phase starts, and the part after. For this reason,
it's sometimes helpful to think of the Gossip handler executing three phases, though the
code between phase 1 and 3 is mostly identical.

From this point, we will consider a single entropy session, initiated by the AntiEntropy
go routine and implemented in the AntiEntropySync method, and handled by the remote peer
using the Gossip handler method. However, it is VERY IMPORTANT to remember that multiple
Gossip handlers can be executing concurrently or that a Gossip and AntiEntropySync can
be executing concurrently on the same peer. Many of the safety decisions made in this
package directly follow from the possibility that two anti-entropy sessions are
happening concurrently.

The connection between the initiator and the remote is a gRPC bidirectional stream. This
means that at any time either the initiator or the remote can send a message to the
other. Messages in one direction are totally ordered; e.g. all messages from the
initiator to the remote are ordered and vice versa. Messages are not ordered between the
two directions. Because of this, the Send and Recv on the stream are independent
operations, meaning that it is possible to concurrently Send and Recv. However, both
Send and Recv are blocking operations, which means it is not possible to concurrently
Send or concurrently Recv. A general rule of thumb is that one and only one go routine
must send messages and one and only one go routine must recv messages.

Each of the two phases in both AntiEntropySync and Gossip are implemented in their own
go routines. All four phases may send a message. Therefore to coordinate the two
initiator phase go routines and the two Gossip phase go routines it is necessary to have
a third go routine that reads messages from a channel and sends those messages one at a
time from each go routine. Only the phase 1 go routines receive messages, so it's not
necessary to have a recv channel. The two phase go routines and sending go routine are
the three go routines that must be synchronized before the anti-entropy session has
concluded successfully.

On the initiator, the phase 1 go routine iterates over all of the objects in the
database and sends CHECK requests as needed to the remote replica. The phase 1 go
routine of the remote replica receives these CHECK requests and compares the version
from the initiator with its local version. If the remote version is earlier than the
initiator version, it sends a CHECK message back to the initiator to retrieve the later
data. If the remote version is later than the initiator version it sends a REPAIR
message back to the initiator with the later data. If the versions are equal, it does
nothing. When the initiator has completed iterating over the database it sends a
COMPLETE message to the remote. When the remote receives the COMPLETE message, it starts
its phase 2 and the phase 1 go routine continues to receive messages from the initiator
with the exception that it will no longer respond to CHECK messages from the initiator.

On the initiator, the phase 2 go routine is started right after its phase 1 go routine.
Its job is to receive CHECK and REPAIR messages from the remote. If it receives a REPAIR
it checks to make sure the repair version is later than the local version, then updates
it. If it receives a CHECK message and the local version is later than the remote's
version it sends the corresponding REPAIR back to the remote. On the remote, the phase 2
go routine iterates over its local database to find objects that were not included in
phase 1 (e.g. objects that are on the remote but not the initiator). It sends REPAIR
messages back to the initiator with those versions. When it is done iterating over its
local database it sends a COMPLETE message back ot the initiator.

The initiator phase 1 and remote phase 2 go routines complete when they've finished
iterating over their database, and both send COMPLETE messages when they're finished.
The initiator phase 2 routine ends when it receives a COMPLETE message from the remote,
at this point it also closes the stream. The remote phase 1 go routine ends when it
receives an EOF, meaning the stream has been closed gracefully. When the recv go
routines (initiator phase 2, remote phase 1) end, they close the sending channel. The
sending go routine ends when it has sent all messages on the channel and the channel is
closed. The anti-entropy session is complete when all go routines are done.

The logic described by these phases and methods is fairly complex since we're dealing
with multiple go routines in two different processes. When working on this code - take
care with synchronization -- this is where most of the bugs and safety issues arise!

At the end of an anti-entropy session both the initiator and the remote will have
identical underlying databases until the next access!
*/
package replica
