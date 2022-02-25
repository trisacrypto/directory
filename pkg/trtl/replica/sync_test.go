package replica_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"github.com/trisacrypto/directory/pkg/trtl/replica"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"google.golang.org/protobuf/proto"
)

func TestSelectZeroOrOnePeer(t *testing.T) {
	fixtures := loadFixtures(t)
	testCases := []struct {
		name     string
		peers    []string
		self     string
		expected string
	}{
		{name: "None", peers: []string{}, self: "raphael", expected: ""},
		{name: "OnlySelf", peers: []string{"raphael"}, self: "raphael", expected: ""},
		{name: "OnlyOther", peers: []string{"michelangelo"}, self: "raphael", expected: "michelangelo"},
		{name: "SelfAndOther", peers: []string{"raphael", "michelangelo"}, self: "raphael", expected: "michelangelo"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbPeers := make([]*peers.Peer, 0, len(tc.peers))
			for _, p := range tc.peers {
				dbPeers = append(dbPeers, fixtures[p])
			}
			db := createDB(t, dbPeers)
			replica := initReplica(t, db, fixtures[tc.self])
			selected := replica.SelectPeer()
			if tc.expected == "" {
				require.Nil(t, selected, "should not have selected a peer")
			} else {
				require.NotNil(t, selected, "should have selected a peer")
				require.Equal(t, fixtures[tc.expected].Id, selected.Id, "selected the wrong peer")
			}
		})
	}
}

// Test that peers are selected uniformly at random.
func TestSelectPeerAll(t *testing.T) {
	fixtures := loadFixtures(t)
	peers := make([]*peers.Peer, 0, len(fixtures))
	for _, p := range fixtures {
		peers = append(peers, p)
	}
	db := createDB(t, peers)
	replica := initReplica(t, db, fixtures["raphael"])

	// Run enough times so that all peers have a chance to be selected.
	counter := make(map[uint64]uint8)
	for i := 0; i < 100; i++ {
		peer := replica.SelectPeer()
		require.NotNil(t, peer, "should have selected a peer")
		counter[peer.Id]++
	}

	// All peers except the local replica should be selected.
	require.Len(t, counter, len(peers)-1, "did not select all peers")
	require.NotContains(t, counter, uint64(fixtures["raphael"].Id), "should not have selected the local replica")

	// The peers should be selected uniformly at random.
	for id, count := range counter {
		if id != uint64(fixtures["raphael"].Id) {
			require.GreaterOrEqual(t, count, uint8(20), "peer was selected too few times")
		}
	}
}

func loadFixtures(t *testing.T) map[string]*peers.Peer {
	// Load peer fixtures
	fixtures := make(map[string]*peers.Peer)
	data, err := ioutil.ReadFile("testdata/peers.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &fixtures), "could not unmarshal peers fixtures")
	require.Len(t, fixtures, 4, "unexpected number of peers fixtures")
	return fixtures
}

func createDB(t *testing.T, fixtures []*peers.Peer) *honu.DB {
	tmp, err := ioutil.TempDir("testdata", "*-db")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmp) })

	db, err := honu.Open("leveldb:///" + tmp)
	require.NoError(t, err)

	// Populate the database with the peers
	for _, p := range fixtures {
		p.Created = time.Now().Format(time.RFC3339)
		p.Modified = p.Created
		data, err := proto.Marshal(p)
		require.NoError(t, err)
		_, err = db.Put([]byte(p.Key()), data, options.WithNamespace(wire.NamespaceReplicas))
		require.NoError(t, err)
	}

	return db
}

func initReplica(t *testing.T, db *honu.DB, self *peers.Peer) *replica.Service {
	// Configure replica
	conf := config.Config{
		Replica: config.ReplicaConfig{
			Enabled:        true,
			PID:            self.Id,
			Name:           self.Name,
			Region:         self.Region,
			GossipInterval: 10 * time.Minute,
			GossipSigma:    1500 * time.Millisecond,
		},
		MTLS: config.MTLSConfig{
			Insecure: true,
		},
	}

	replica, err := replica.New(conf, db, nil)
	require.NoError(t, err)
	return replica
}
