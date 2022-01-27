/*
Package replica implements bi-lateral anti-entropy for the trtl database.

TODO: finish documentation.

Gossip implements bilateral anti-entropy: during a Gossip session the initiating
replica pushes updates to the remote peer and pulls requested changes. Using
bidirectional streaming, the initiating peer sends data-less sync messages with
the versions of objects it stores locally. The remote replica then responds with
data if its local version is later or sends a sync message back requesting the
data from the initiating replica if its local version is earlier. No exchange
occurs if both replicas have the same version. At the end of a gossip session,
both replicas should have synchronized and have identical underlying data stores.
*/
package replica
