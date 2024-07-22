package trtl_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// Config values for test replica
const (
	benchRegion = "california"
	benchOwner  = "arnold"
	benchPID    = 42
	benchbuf    = 1024 * 1024
)

type trtlBench struct {
	trtl  *trtl.Server
	conf  config.Config
	grpc  *bufconn.GRPCListener
	tmpdb string
}

// setupTrtl needs to load fixtures, run a Trtl server, and perform cleanup
// similar to all the work done by the trtlTestSuite in server_test.go
// Unfortunately testify test suite does not work with benchmark, so much of
// this work is duplicated here.
func setupTrtl(t testing.TB) (bench *trtlBench, err error) {

	bench = new(trtlBench)

	// Create a tmp directory for the database
	if bench.tmpdb, err = os.MkdirTemp("", "trtldb-*"); err != nil {
		return nil, fmt.Errorf("could not create tmpdb: %w", err)
	}

	// Manually create a configuration
	c := config.Config{
		Maintenance: false,
		BindAddr:    ":4442",
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  true,
		Database: config.DatabaseConfig{
			URL:           fmt.Sprintf("leveldb:///%s", bench.tmpdb),
			ReindexOnBoot: false,
		},
		Replica: config.ReplicaConfig{
			Enabled: false, // Replica is tested in the replica package
			PID:     benchPID,
			Region:  benchRegion,
			Name:    benchOwner,
		},
		MTLS: config.MTLSConfig{
			Insecure: true,
		},
		Backup: config.BackupConfig{
			Enabled: false,
		},
	}

	// Mark as processed since the config wasn't loaded from the environment
	if bench.conf, err = c.Mark(); err != nil {
		return nil, err
	}

	// Create the trtl server
	if bench.trtl, err = trtl.New(bench.conf); err != nil {
		return nil, err
	}

	// Create a bufconn listener(s) so that there are no actual network requests
	bench.grpc = bufconn.New("")

	// Run the test server without signals, background routines or maintenance mode checks
	go bench.trtl.Run(bench.grpc.Listener)
	t.Cleanup(func() { cleanup(bench) })

	return bench, err
}

// cleanup the current temporary directory, configuration, and running services.
func cleanup(bench *trtlBench) (err error) {
	// Shutdown the trtl server if it is running
	// This should shutdown all the running services and close the database
	// Note that Shutdown should be graceful and not shutdown anything not running.
	if bench.trtl != nil {
		if err = bench.trtl.Shutdown(); err != nil {
			return err
		}
	}

	// Shutdown the gRPC connection if it's running
	if bench.grpc != nil {
		bench.grpc.Release()
	}

	// Cleanup the benchdb and delete any stray files
	if bench.tmpdb != "" {
		os.RemoveAll(bench.tmpdb)
	}

	return nil
}

func BenchmarkTrtlGet(b *testing.B) {
	// Run setupTrtl
	bench, e := setupTrtl(b)
	require.NotNil(b, bench)
	require.NoError(b, e)

	// Start the gRPC client.
	require.NoError(b, bench.grpc.Connect(context.Background()))
	defer bench.grpc.Close()
	cc := pb.NewTrtlClient(bench.grpc.Conn)
	ctx := context.Background()

	// Manually add some fixtures
	// Create a put request
	key := []byte("terminator")
	catchphrase := []byte("hasta la vista, baby")
	_, err := cc.Put(ctx, &pb.PutRequest{
		Key:   key,
		Value: catchphrase,
	})
	require.NoError(b, err)

	// Run the trtl Get on a loop
	b.ResetTimer()
	var gErr error
	for i := 0; i < b.N; i++ {
		_, gErr = cc.Get(ctx, &pb.GetRequest{Key: key})
	}
	require.NoError(b, gErr)
}

func BenchmarkTrtlPut(b *testing.B) {
	// Run setupTrtl
	bench, e := setupTrtl(b)
	require.NotNil(b, bench)
	require.NoError(b, e)

	// Start the gRPC client.
	require.NoError(b, bench.grpc.Connect(context.Background()))
	defer bench.grpc.Close()
	cc := pb.NewTrtlClient(bench.grpc.Conn)
	ctx := context.Background()

	// Manually create some fixtures
	key := []byte("sarah")
	catchphrase := []byte("How's the knee?")

	// Run the trtl Get on a loop
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		_, err = cc.Put(ctx, &pb.PutRequest{
			Key:   key,
			Value: catchphrase,
		})
	}
	require.NoError(b, err)

}

func BenchmarkTrtlDelete(b *testing.B) {
	// Run setupTrtl
	bench, e := setupTrtl(b)
	require.NotNil(b, bench)
	require.NoError(b, e)

	// Start the gRPC client.
	require.NoError(b, bench.grpc.Connect(context.Background()))
	defer bench.grpc.Close()
	cc := pb.NewTrtlClient(bench.grpc.Conn)
	ctx := context.Background()

	// Reset the timer to focus only on the Delete
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		// Stop timer and put key/value to delete
		b.StopTimer()
		_, err = cc.Put(ctx, &pb.PutRequest{
			Key:   []byte(fmt.Sprintf("t%d", i)),
			Value: []byte(strconv.Itoa(i)),
		})
		require.NoError(b, err)

		// Start timer and perform the Delete
		b.StartTimer()
		_, err = cc.Delete(ctx, &pb.DeleteRequest{
			Key: []byte(fmt.Sprintf("t%d", i)),
		})
	}
	require.NoError(b, err)
}

func BenchmarkTrtlIter(b *testing.B) {
	// Run setupTrtl
	bench, e := setupTrtl(b)
	require.NotNil(b, bench)
	require.NoError(b, e)

	// Start the gRPC client.
	require.NoError(b, bench.grpc.Connect(context.Background()))
	defer bench.grpc.Close()
	cc := pb.NewTrtlClient(bench.grpc.Conn)
	ctx := context.Background()

	// Manually add some fixtures
	var err error
	for i := 0; i < 1000; i++ {
		_, err = cc.Put(ctx, &pb.PutRequest{
			Namespace: "terminators",
			Key:       []byte(fmt.Sprintf("t%d", i)),
			Value:     []byte(strconv.Itoa(i)),
		})
	}
	require.NoError(b, err)

	// Run the trtl Iter on a loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = cc.Iter(ctx, &pb.IterRequest{
			Namespace: "terminators",
			Prefix:    []byte("t"),
		})
	}
	require.NoError(b, err)
}

func BenchmarkTrtlCursor(b *testing.B) {
	// Run setupTrtl
	bench, e := setupTrtl(b)
	require.NotNil(b, bench)
	require.NoError(b, e)

	// Start the gRPC client.
	require.NoError(b, bench.grpc.Connect(context.Background()))
	defer bench.grpc.Close()
	cc := pb.NewTrtlClient(bench.grpc.Conn)
	ctx := context.Background()

	// Manually add some fixtures
	var err error
	for i := 0; i < 1000; i++ {
		_, err = cc.Put(ctx, &pb.PutRequest{
			Namespace: "terminators",
			Key:       []byte(fmt.Sprintf("t%d", i)),
			Value:     []byte(strconv.Itoa(i)),
		})
	}
	require.NoError(b, err)

	var stream pb.Trtl_CursorClient
	stream, err = cc.Cursor(ctx, &pb.CursorRequest{Namespace: "terminators"})
	require.NoError(b, err, "could not create cursor stream")

	// Run the trtl Cursor on a loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
		}
		require.NoError(b, err)
	}
	require.NoError(b, err)
}

func BenchmarkTrtlBatch(b *testing.B) {
	// Run setupTrtl
	bench, e := setupTrtl(b)
	require.NotNil(b, bench)
	require.NoError(b, e)

	// Start the gRPC client.
	require.NoError(b, bench.grpc.Connect(context.Background()))
	defer bench.grpc.Close()
	cc := pb.NewTrtlClient(bench.grpc.Conn)
	ctx := context.Background()

	// Manually add some fixtures
	var err error
	requests := map[int64]*pb.BatchRequest{
		1: {
			Id: 1,
			Request: &pb.BatchRequest_Put{
				Put: &pb.PutRequest{
					Key:   []byte("john"),
					Value: []byte("easy money"),
				},
			},
		},
		2: {
			Id: 2,
			Request: &pb.BatchRequest_Put{
				Put: &pb.PutRequest{
					Key:   []byte("terminator"),
					Value: []byte("i'll be back"),
				},
			},
		},
	}
	stream, err := cc.Batch(ctx)
	require.NoError(b, err)

	// Run the trtl Batch on a loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range requests {
			err = stream.Send(r)
			require.NoError(b, err)
		}
	}
	_, err = stream.CloseAndRecv()
	require.NoError(b, err)
}
