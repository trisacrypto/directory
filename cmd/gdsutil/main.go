package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "gdsutil"
	app.Version = pkg.Version()
	app.Usage = "utilities for operating the GDS service and database"
	app.Commands = []cli.Command{
		{
			Name:     "ldb:keys",
			Usage:    "list the keys currently in the leveldb store",
			Category: "leveldb",
			Action:   ldbKeys,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "s, stringify",
					Usage: "stringify keys otherwise they are base64 encoded",
				},
				cli.StringFlag{
					Name:  "p, prefix",
					Usage: "specify a prefix to filter keys on",
				},
			},
		},
		{
			Name:      "ldb:get",
			Usage:     "get the value for the specified key",
			ArgsUsage: "key [key ...]",
			Category:  "leveldb",
			Action:    ldbGet,
			Before:    openLevelDB,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "b, b64decode",
					Usage: "specify the keys as base64 encoded values which must be decoded",
				},
				cli.StringFlag{
					Name:  "o, out",
					Usage: "write the fetched key to directory if specified, otherwise printed",
				},
			},
		},
		{
			Name:      "ldb:put",
			Usage:     "insert key/value pair into database, loading from disk if necessary",
			ArgsUsage: "key/path [value]\n\n   Examples:\n   gdsutil ldb:put foo bar\n   gdsutil ldb:put foo bar.json\n   gdsutil ldb:put foo.json",
			Category:  "leveldb",
			Action:    ldbPut,
			Before:    openLevelDB,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "b, b64decode",
					Usage: "specify the key and value as base64 encoded strings which must be decoded",
				},
				cli.StringFlag{
					Name:  "f, format",
					Usage: "format of the data (raw, json, pb)",
					Value: "json",
				},
			},
		},
		{
			Name:      "ldb:delete",
			Usage:     "delete the leveldb record for the specified key(s)",
			ArgsUsage: "key [key ...]",
			Category:  "leveldb",
			Action:    ldbDelete,
			Before:    openLevelDB,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "b, b64decode",
					Usage: "specify the keys as base64 encoded values which must be decoded",
				},
			},
		},
		{
			Name:      "peers:add",
			Usage:     "add peers to the network by pid",
			ArgsUsage: "pid",
			Category:  "replica",
			Before:    openLevelDB,
			Action:    addPeers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				// TODO allow the user to add multiple peers at a time?
				cli.Uint64Flag{
					Name:  "p, pid",
					Usage: "specify the pid for the peer to add",
				},
			},
		},
		{
			Name:      "peers:delete",
			Usage:     "tombstone a peer by pid (does not remove from ldb)",
			ArgsUsage: "pid",
			Category:  "replica",
			Before:    openLevelDB,
			Action:    delPeers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				// TODO allow the user to rm multiple peers at a time?
				cli.Uint64Flag{
					Name:  "p, pid",
					Usage: "specify the pid for the peer to tombstone",
				},
			},
		},
		{
			Name:     "peers:list",
			Usage:    "get a status report of all peers in the network",
			Category: "replica",
			Before:   openLevelDB,
			Action:   listPeers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				// TODO: have we standardized on how to reference regions?
				cli.StringSliceFlag{
					Name:  "r, region",
					Usage: "specify a region for peers to be returned",
				},
				cli.BoolFlag{
					Name:  "s, status",
					Usage: "specify for status-only, will not return peer details",
				},
			},
		},
		{
			Name:      "gossip",
			Usage:     "initiate a gossip session with a remote replica (for debugging)",
			ArgsUsage: "remote:port",
			Category:  "replica",
			Before:    openLevelDB,
			Action:    gossip,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "p, partial",
					Usage: "ignore any objects not specified in request",
				},
				cli.StringSliceFlag{
					Name:  "n, namespaces",
					Usage: "specify the namespaces to replicate (if empty, all are replicated)",
				},
				cli.StringSliceFlag{
					Name:  "o, objects",
					Usage: "specify the object keys to replicate (otherwise all objects from namespaces will be used)",
				},
				cli.BoolFlag{
					Name:  "D, dryrun",
					Usage: "show changes that would occur, does not modify database",
				},
			},
		},
		{
			Name:     "gossip:migrate",
			Usage:    "migrate objects to replication context",
			Category: "replica",
			Before:   openLevelDB,
			Action:   gossipMigrate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.StringFlag{
					Name:   "a, addr",
					Usage:  "bind addr of the local replica (for name processing)",
					EnvVar: "GDS_REPLICA_BIND_ADDR",
				},
				cli.Uint64Flag{
					Name:   "p, pid",
					Usage:  "process id of the local replica",
					EnvVar: "GDS_REPLICA_PID",
				},
				cli.StringFlag{
					Name:   "r, region",
					Usage:  "geographic region of the local replica",
					EnvVar: "GDS_REPLICA_REGION",
				},
				cli.StringFlag{
					Name:   "n, name",
					Usage:  "human readable name of the local replica",
					EnvVar: "GDS_REPLICA_NAME",
				},
				cli.BoolFlag{
					Name:  "D, dryrun",
					Usage: "show changes that would occur, does not modify database",
				},
			},
		},
		{
			Name:      "decrypt",
			Usage:     "decrypt base64 encoded ciphertext with an HMAC signature",
			ArgsUsage: "ciphertext hmac",
			Category:  "cipher",
			Action:    cipherDecrypt,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "k, key",
					Usage:  "secret key to decrypt the cipher text",
					EnvVar: "GDS_SECRET_KEY",
				},
			},
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// LevelDB Actions
//===========================================================================

var (
	ldb            *leveldb.DB
	vaspPrefix     = []byte("vasps::")
	certreqsPrefix = []byte("certreqs::")
	// replicaPrefix  = []byte("peers::")
	indexPrefix    = []byte("index::")
	sequencePrefix = []byte("sequence::")
	indexNames     = []byte("index::names")
	indexCountries = []byte("index::countries")
	sequencePK     = []byte("sequence::pks")
)

func ldbKeys(c *cli.Context) (err error) {
	defer ldb.Close()

	var prefix *util.Range
	if prefixs := c.String("prefix"); prefixs != "" {
		prefix = util.BytesPrefix([]byte(prefixs))
	}

	iter := ldb.NewIterator(prefix, nil)
	defer iter.Release()

	stringify := c.Bool("stringify")
	for iter.Next() {
		if stringify {
			fmt.Printf("- %s\n", string(iter.Key()))
		} else {
			fmt.Printf("- %s\n", base64.RawStdEncoding.EncodeToString(iter.Key()))
		}
	}

	if err = iter.Error(); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func ldbGet(c *cli.Context) (err error) {
	defer ldb.Close()

	var out string
	if out = c.String("out"); out != "" {
		// Check that out is a directory
		var info fs.FileInfo
		if info, err = os.Stat(out); err != nil {
			return cli.NewExitError(err, 1)
		}
		if !info.IsDir() {
			return cli.NewExitError("specify a directory to write files out to", 1)
		}
	}

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args() {
		var key []byte
		if b64decode {
			if key, err = base64.RawStdEncoding.DecodeString(keys); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			key = []byte(keys)
		}

		var data []byte
		if data, err = ldb.Get(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}

		// Unmarshal the thing
		var (
			jsonValue interface{}
			pbValue   proto.Message
		)

		// Determine how to unmarshal the data
		if bytes.HasPrefix(key, vaspPrefix) {
			vasp := new(pb.VASP)
			if err = proto.Unmarshal(data, vasp); err != nil {
				return cli.NewExitError(err, 1)
			}
			pbValue = vasp
		} else if bytes.HasPrefix(key, certreqsPrefix) {
			careq := new(models.CertificateRequest)
			if err = proto.Unmarshal(data, careq); err != nil {
				return cli.NewExitError(err, 1)
			}
			pbValue = careq
		} else if bytes.Equal(key, indexNames) {
			jsonValue = make(map[string]string)
			if err = unmarshalGZJSON(data, &jsonValue); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, indexCountries) {
			jsonValue = make(map[string][]string)
			if err = unmarshalGZJSON(data, &jsonValue); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.HasPrefix(key, sequencePrefix) {
			pk, n := binary.Uvarint(data)
			if n <= 0 {
				return cli.NewExitError("could not parse sequence", 1)
			}
			jsonValue = pk
		} else {
			return cli.NewExitError("could not determine unmarshal type", 1)
		}

		// Marshal JSON representation
		var outdata []byte
		if jsonValue != nil {
			if outdata, err = json.MarshalIndent(jsonValue, "", "  "); err != nil {
				return cli.NewExitError(err, 1)
			}
		}
		if pbValue != nil {
			jsonpb := protojson.MarshalOptions{
				Multiline:       true,
				Indent:          "  ",
				AllowPartial:    true,
				UseProtoNames:   true,
				UseEnumNumbers:  false,
				EmitUnpopulated: true,
			}
			if outdata, err = jsonpb.Marshal(pbValue); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if out != "" {
			path := filepath.Join(out, string(key)+".json")
			if err = ioutil.WriteFile(path, outdata, 0644); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			fmt.Println(string(outdata) + "\n")
		}
	}

	return nil
}

func ldbPut(c *cli.Context) (err error) {
	defer ldb.Close()

	if c.NArg() == 0 || c.NArg() > 2 {
		return cli.NewExitError("specify path, key and path, or key and value as arguments", 1)
	}

	// Determine the key and value as follows:
	// If only one argument, assume it is a path; key is basename, value is data
	// If two arguments, check if second value is a path, otherwise second value is a key
	args := c.Args()
	format := strings.ToLower(c.String("format"))
	var key, data, value []byte
	if c.NArg() == 1 {
		path := args.Get(0)
		name := filepath.Base(path)
		ext := filepath.Ext(name)
		if strings.TrimLeft(ext, ".") != format {
			return cli.NewExitError(fmt.Errorf("mismatch file extension %q and data format %q: specify --format", ext, format), 1)
		}

		key = []byte(strings.TrimSuffix(name, ext))
		if data, err = ioutil.ReadFile(path); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		key = []byte(args.Get(0))

		// Determine if second argument is a path
		varg := args.Get(1)
		if isFile(varg) {
			ext := filepath.Ext(varg)
			if strings.TrimLeft(ext, ".") != format {
				return cli.NewExitError(fmt.Errorf("mismatch file extension %q and data format %q: specify --format", ext, format), 1)
			}
			if data, err = ioutil.ReadFile(varg); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			data = []byte(varg)
		}

	}

	// Perform base64 decoding if necessary
	b64decode := c.Bool("b64decode")
	if b64decode {
		if key, err = base64.RawStdEncoding.DecodeString(string(key)); err != nil {
			return cli.NewExitError(err, 1)
		}
		if data, err = base64.RawStdEncoding.DecodeString(string(data)); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	// Quick spot check
	if len(data) == 0 || len(key) == 0 {
		return cli.NewExitError("no key or value found", 1)
	}

	// Unmarshal the thing from data then
	// Marshal the database representation
	if format != "raw" && format != "bytes" {
		jsonpb := &protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		if bytes.HasPrefix(key, vaspPrefix) {
			vasp := new(pb.VASP)
			switch format {
			case "json":
				if err = jsonpb.Unmarshal(data, vasp); err != nil {
					return cli.NewExitError(err, 1)
				}
				if value, err = proto.Marshal(vasp); err != nil {
					return cli.NewExitError(err, 1)
				}
			case "pb", "proto", "protobuf":
				if err = proto.Unmarshal(data, vasp); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = data
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal VASP format %q", format), 1)
			}

		} else if bytes.HasPrefix(key, certreqsPrefix) {
			careq := new(models.CertificateRequest)
			switch format {
			case "json":
				if err = jsonpb.Unmarshal(data, careq); err != nil {
					return cli.NewExitError(err, 1)
				}
				if value, err = proto.Marshal(careq); err != nil {
					return cli.NewExitError(err, 1)
				}
			case "pb", "proto", "protobuf":
				if err = proto.Unmarshal(data, careq); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = data
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal Certificate Request format %q", format), 1)
			}
		} else if bytes.HasPrefix(key, indexPrefix) {
			switch format {
			case "json":
				// gzip compress the index data
				buf := &bytes.Buffer{}
				gz := gzip.NewWriter(buf)
				if _, err = gz.Write(data); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = buf.Bytes()
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal index format %q", format), 1)
			}
		} else if bytes.Equal(key, sequencePK) {
			switch format {
			case "json":
				var pk uint64
				if err = json.Unmarshal(data, &pk); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = make([]byte, binary.MaxVarintLen64)
				binary.PutUvarint(value, pk)
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal sequence format %q", format), 1)
			}

		} else {
			return cli.NewExitError("could not determine unmarshal type from key", 1)
		}
	} else {
		// Raw or bytes data is just the data
		value = data
	}

	// Final spot check
	if len(value) == 0 {
		return cli.NewExitError("no value marshaled", 1)
	}

	// Put the key/value to the database
	if err = ldb.Put(key, value, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func ldbDelete(c *cli.Context) (err error) {
	defer ldb.Close()
	if c.NArg() == 0 {
		return cli.NewExitError("specify at least one key to delete", 1)
	}

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args() {
		var key []byte
		if b64decode {
			if key, err = base64.RawStdEncoding.DecodeString(keys); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			key = []byte(keys)
		}

		if err = ldb.Delete(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

//===========================================================================
// LevelDB Helper Functions
//===========================================================================

func openLevelDB(c *cli.Context) (err error) {
	var uri string
	if uri = c.String("db"); uri == "" {
		return cli.NewExitError("specify path to leveldb database", 1)
	}

	var dsn *store.DSN
	if dsn, err = store.ParseDSN(uri); err != nil {
		return cli.NewExitError(err, 1)
	}

	if dsn.Scheme != "leveldb" && dsn.Scheme != "ldb" {
		return cli.NewExitError("this action requires a leveldb DSN", 1)
	}

	if ldb, err = leveldb.OpenFile(dsn.Path, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return fi.Mode().IsRegular()
}

func unmarshalGZJSON(data []byte, val interface{}) (err error) {
	buf := bytes.NewBuffer(data)
	var gz *gzip.Reader
	if gz, err = gzip.NewReader(buf); err != nil {
		return fmt.Errorf("could not decompress data: %s", err)
	}
	decoder := json.NewDecoder(gz)
	if err = decoder.Decode(&val); err != nil {
		return fmt.Errorf("could not decode json data: %s", err)
	}
	return nil
}

//lint:ignore U1000 leaving this function to pair with unmarshalGZJSON in the future
func marshalGZJSON(val interface{}) (data []byte, err error) {
	var buf *bytes.Buffer
	gz := gzip.NewWriter(buf)
	encoder := json.NewEncoder(gz)
	if err = encoder.Encode(val); err != nil {
		return data, fmt.Errorf("could not encode data: %s", err)
	}
	if err = gz.Close(); err != nil {
		return data, fmt.Errorf("could not compress data: %s", err)
	}
	return buf.Bytes(), nil
}

//===========================================================================
// Peer Management Replica Actions
//===========================================================================

func addPeers(c *cli.Context) (err error) {
	defer ldb.Close()
	if c.NArg() != 1 {
		return cli.NewExitError("must specify pid for peer", 1)
	}
	pid := c.Uint64("pid")
	peer := &peers.Peer{
		Id: pid,
	}

	// initialize a client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.Args()[0], grpc.WithInsecure()); err != nil {
		return cli.NewExitError(err, 1)
	}
	client := peers.NewPeerManagementClient(cc)

	// create a new context and pass the parent context in
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// call client.AddPeer with the pid
	var out *peers.PeersStatus
	if out, err = client.AddPeers(ctx, peer); err != nil {
		return cli.NewExitError(err, 1)
	}
	// print the status
	printJSON(out)

	return nil
}

func delPeers(c *cli.Context) (err error) {
	defer ldb.Close()
	if c.NArg() != 1 {
		return cli.NewExitError("must specify pid for peer", 1)
	}
	pid := c.Uint64("pid")
	peer := &peers.Peer{
		Id: pid,
	}

	// initialize a client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.Args()[0], grpc.WithInsecure()); err != nil {
		return cli.NewExitError(err, 1)
	}
	client := peers.NewPeerManagementClient(cc)

	// create a new context and pass the parent context in
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// call client.RmPeer with the pid
	var out *peers.PeersStatus
	if out, err = client.RmPeers(ctx, peer); err != nil {
		return cli.NewExitError(err, 1)
	}

	// print the status
	printJSON(out)

	return nil
}

func listPeers(c *cli.Context) (err error) {
	defer ldb.Close()

	// determine if this is a region-specific or status only request
	so := c.GlobalBool("status")
	regions := c.StringSlice("region")
	filter := &peers.PeersFilter{
		Region:     regions,
		StatusOnly: so,
	}

	// initialize a client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.Args()[0], grpc.WithInsecure()); err != nil {
		return cli.NewExitError(err, 1)
	}
	client := peers.NewPeerManagementClient(cc)

	// create a new context and pass the parent context in
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// call client.GetPeers with filter
	var out *peers.PeersList
	if out, err = client.GetPeers(ctx, filter); err != nil {
		return cli.NewExitError(err, 1)
	}
	// print the peers
	printJSON(out)

	return nil
}

//===========================================================================
// Replica Gossip Actions
//===========================================================================

func gossip(c *cli.Context) (err error) {
	defer ldb.Close()
	if c.NArg() != 1 {
		return cli.NewExitError("must specify the remote replica addr:port", 1)
	}

	// If objects is specified then load them from the database
	objs := c.StringSlice("objects")
	namespaces := c.StringSlice("namespaces")

	versions := &global.VersionVectors{
		Objects:    make([]*global.Object, 0),
		Partial:    c.Bool("partial"),
		Namespaces: namespaces,
	}

	// Sanity check to ensure no duplicates and no ignored objects
	if len(objs) > 0 && len(namespaces) > 0 {
		return cli.NewExitError("specify objects or namespaces not both", 1)
	}

	// Load objects
	if len(objs) > 0 {
		for _, key := range objs {
			var obj *global.Object
			if obj, err = loadMetadata(key); err != nil {
				return cli.NewExitError(err, 1)
			}
			versions.Objects = append(versions.Objects, obj)
		}
	} else {

		if len(versions.Namespaces) == 0 {
			// Specify "all" namespaces manually (opt-in to what is replicated).
			// NOTE: if there is a namespace being omitted from gossip without being
			// specified in the command line, it's likely it needs to be added here.
			namespaces = store.Namespaces[:]
		}

		for _, ns := range namespaces {
			var objs []*global.Object
			if objs, err = loadNamespaceMetadata(ns); err != nil {
				return cli.NewExitError(err, 1)
			}
			versions.Objects = append(versions.Objects, objs...)
		}
	}

	// Initialize the gossip client
	// TODO: move to its own helper function
	// TODO: make secure
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.Args()[0], grpc.WithInsecure()); err != nil {
		return cli.NewExitError(err, 1)
	}
	client := global.NewReplicationClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	var rep *global.Updates
	if rep, err = client.Gossip(ctx, versions); err != nil {
		return cli.NewExitError(err, 1)
	}

	// If dryrun only print objects retrieved from remote; do not modify database.
	if c.Bool("dryrun") {
		fmt.Printf("sent %d versions to remote\n", len(versions.Objects))
		fmt.Printf("received %d repairs from remote\n", len(rep.Objects))
		for _, obj := range rep.Objects {
			fmt.Printf("  - %s (v%d.p%d)\n", obj.Key, obj.Version.Version, obj.Version.Pid)
		}
		return nil
	}

	// Repair local database as last step if not dryrun
	for _, obj := range rep.Objects {
		// Get the struct from the any in the object
		// This isn't totally necessary but is some sanity checking to make sure things aren't totally out of whack
		var msg proto.Message
		if msg, err = obj.Data.UnmarshalNew(); err != nil {
			return cli.NewExitError(fmt.Errorf("could not unmarshal %s in %s from any: %s", obj.Key, obj.Namespace, err), 1)
		}

		switch msg := msg.(type) {
		case *pb.VASP:
			if obj.Namespace != store.NamespaceVASPs {
				return cli.NewExitError(fmt.Errorf("type/namespace mismatch %s in %s from any: %T", obj.Key, obj.Namespace, msg), 1)
			}
		case *models.CertificateRequest:
			if obj.Namespace != store.NamespaceCertReqs {
				return cli.NewExitError(fmt.Errorf("type/namespace mismatch %s in %s from any: %T", obj.Key, obj.Namespace, msg), 1)
			}
		default:
			return cli.NewExitError(fmt.Errorf("could not handle %s in %s from any type %T", obj.Key, obj.Namespace, msg), 1)
		}

		// Store the data in the database
		if err = ldb.Put([]byte(obj.Key), obj.Data.Value, nil); err != nil {
			return cli.NewExitError(fmt.Errorf("could not put %s in %s to database: %s", obj.Key, obj.Namespace, err), 1)
		}
	}

	fmt.Printf("sent %d versions to remote\n", len(versions.Objects))
	fmt.Printf("received %d repairs from remote\n", len(rep.Objects))
	return nil
}

func gossipMigrate(c *cli.Context) (err error) {
	// Create a replica config to create a new version manager
	conf := config.ReplicaConfig{
		Enabled:  true,
		BindAddr: c.String("addr"),
		PID:      c.Uint64("pid"),
		Region:   c.String("region"),
		Name:     c.String("name"),
	}
	dryrun := c.Bool("dryrun")

	if err = conf.Validate(); err != nil {
		return cli.NewExitError(err, 1)
	}

	var vm *global.VersionManager
	if vm, err = global.New(conf); err != nil {
		return cli.NewExitError(err, 1)
	}

	// Migrate VASP objects
	migrated := 0
	iter := ldb.NewIterator(util.BytesPrefix([]byte("vasps::")), nil)
vaspLoop:
	for iter.Next() {
		vasp := &pb.VASP{}
		if err = proto.Unmarshal(iter.Value(), vasp); err != nil {
			fmt.Printf("could not unmarshal %q: %s\n", iter.Key(), err)
			continue vaspLoop
		}

		key := string(iter.Key())
		meta, deletedOn, err := models.GetMetadata(vasp)
		if err != nil {
			fmt.Printf("could not get metadata %q: %s\n", key, err)
			continue vaspLoop
		}

		if !deletedOn.IsZero() {
			// Handle tombstone
			fmt.Printf("cannot handle tombstone %q\n", key)
			continue vaspLoop
		}

		if meta == nil {
			meta = &global.Object{}
		}

		// Update metadata
		meta.Key = key
		meta.Namespace = store.NamespaceVASPs
		meta.Owner = vm.Owner
		meta.Region = vm.Region
		meta.Version = &global.Version{
			Pid:     vasp.Version.Pid,
			Version: vasp.Version.Version,
		}
		if err = vm.Update(meta); err != nil {
			fmt.Printf("could not update metadata %q: %s\n", iter.Key(), err)
			continue vaspLoop
		}

		if err = models.SetMetadata(vasp, meta, deletedOn); err != nil {
			fmt.Printf("could not set metadata %q: %s\n", key, err)
			continue vaspLoop
		}

		if !dryrun {
			var data []byte
			if data, err = proto.Marshal(vasp); err != nil {
				fmt.Printf("could not marshal %q: %s\n", iter.Key(), err)
				continue vaspLoop
			}

			if err = ldb.Put([]byte(key), data, nil); err != nil {
				fmt.Printf("could not write %q: %s\n", iter.Key(), err)
				continue vaspLoop
			}
		} else {
			fmt.Printf("%+v\n", meta)
		}

		migrated++
	}

	if err = iter.Error(); err != nil {
		iter.Release()
		return cli.NewExitError(fmt.Errorf("could not iterate over VASPs: %s", err), 1)
	}
	iter.Release()

	fmt.Printf("migrated %d VASP objects\n", migrated)

	// Migrate CertReq objects
	migrated = 0
	iter = ldb.NewIterator(util.BytesPrefix([]byte("certreqs::")), nil)
certreqLoop:
	for iter.Next() {
		certreq := &models.CertificateRequest{}
		if err = proto.Unmarshal(iter.Value(), certreq); err != nil {
			fmt.Printf("could not unmarshal %q: %s\n", iter.Key(), err)
			continue certreqLoop
		}

		key := string(iter.Key())
		if certreq.Deleted != "" {
			// Handle tombstone
			fmt.Printf("cannot handle tombstone %q\n", key)
			continue certreqLoop
		}

		// Create metadata if it does not exist
		if certreq.Metadata == nil {
			certreq.Metadata = &global.Object{}
		}

		// Update metadata
		certreq.Metadata.Key = key
		certreq.Metadata.Namespace = store.NamespaceCertReqs
		certreq.Metadata.Owner = vm.Owner
		certreq.Metadata.Region = vm.Region
		certreq.Metadata.Version = &global.Version{
			Pid:     0,
			Version: 2,
		}
		if err = vm.Update(certreq.Metadata); err != nil {
			fmt.Printf("could not update metadata %q: %s\n", iter.Key(), err)
			continue certreqLoop
		}

		if !dryrun {
			var data []byte
			if data, err = proto.Marshal(certreq); err != nil {
				fmt.Printf("could not marshal %q: %s\n", iter.Key(), err)
				continue certreqLoop
			}

			if err = ldb.Put([]byte(key), data, nil); err != nil {
				fmt.Printf("could not write %q: %s\n", iter.Key(), err)
				continue certreqLoop
			}
		} else {
			fmt.Printf("%+v\n", certreq.Metadata)
		}

		migrated++
	}

	if err = iter.Error(); err != nil {
		iter.Release()
		return cli.NewExitError(fmt.Errorf("could not iterate over CertificateRequests: %s", err), 1)
	}
	iter.Release()

	fmt.Printf("migrated %d CertificateRequest objects\n", migrated)
	return nil
}

// Helper function to load object metadata from leveldb
func loadMetadata(key string) (obj *global.Object, err error) {
	// Load object from the data
	var data []byte
	if data, err = ldb.Get([]byte(key), nil); err != nil {
		return nil, fmt.Errorf("could not get %q: %s", key, err)
	}

	// Detect the type of object, deserialize, and extract object metadata
	namespace := strings.Split(key, ":")[0]
	switch namespace {
	case store.NamespaceVASPs:
		vasp := &pb.VASP{}
		if err = proto.Unmarshal(data, vasp); err != nil {
			return nil, fmt.Errorf("could not unmarshal %q into vasp: %s", key, err)
		}
		obj, _, err = models.GetMetadata(vasp)
		return obj, err
	case store.NamespaceCertReqs:
		careq := &models.CertificateRequest{}
		if err = proto.Unmarshal(data, careq); err != nil {
			return nil, fmt.Errorf("could not unmarshal %q into certreq: %s", key, err)
		}
		return careq.Metadata, nil
	case "peers":
		peer := &peers.Peer{}
		if err = proto.Unmarshal(data, peer); err != nil {
			return nil, fmt.Errorf("could not unmarshal %q into peer: %s", key, err)
		}
		return peer.Metadata, nil
	default:
		return nil, fmt.Errorf("could not parse namespace %q", namespace)
	}
}

// Helper function to load all object metadata for a namespace
func loadNamespaceMetadata(ns string) (objs []*global.Object, err error) {
	prefix := util.BytesPrefix([]byte(ns))
	iter := ldb.NewIterator(prefix, nil)
	defer iter.Release()

	objs = make([]*global.Object, 0)
	for iter.Next() {
		var obj *global.Object
		if obj, err = loadMetadata(string(iter.Key())); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
}

//===========================================================================
// Cipher Actions
//===========================================================================

const nonceSize = 12

func cipherDecrypt(c *cli.Context) (err error) {
	if c.NArg() != 2 {
		return cli.NewExitError("must specify ciphertext and hmac arguments", 1)
	}

	var secret string
	if secret = c.String("key"); secret == "" {
		return cli.NewExitError("cipher key required", 1)
	}

	var ciphertext, signature []byte
	if ciphertext, err = base64.RawStdEncoding.DecodeString(c.Args()[0]); err != nil {
		return cli.NewExitError(fmt.Errorf("could not decode ciphertext: %s", err), 1)
	}
	if signature, err = base64.RawStdEncoding.DecodeString(c.Args()[1]); err != nil {
		return cli.NewExitError(fmt.Errorf("could not decode signature: %s", err), 1)
	}

	if len(ciphertext) == 0 {
		return cli.NewExitError("empty cipher text", 1)
	}

	// Create a 32 byte signature of the key
	hash := sha256.New()
	hash.Write([]byte(secret))
	key := hash.Sum(nil)

	// Separate the data from the nonce
	data := ciphertext[:len(ciphertext)-nonceSize]
	nonce := ciphertext[len(ciphertext)-nonceSize:]

	// Validate HMAC signature
	if err = validateHMAC(key, data, signature); err != nil {
		return cli.NewExitError(err, 1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	plainbytes, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(string(plainbytes))
	return nil
}

//===========================================================================
// Cipher Helper Functions
//===========================================================================

func createHMAC(key, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot sign empty data")
	}
	hm := hmac.New(sha256.New, key)
	hm.Write(data)
	return hm.Sum(nil), nil
}

func validateHMAC(key, data, sig []byte) error {
	hmac, err := createHMAC(key, data)
	if err != nil {
		return err
	}

	if !bytes.Equal(sig, hmac) {
		return errors.New("HMAC mismatch")
	}
	return nil
}

//===========================================================================
// Helper Functions
//===========================================================================

// helper function to print JSON response and exit
func printJSON(m proto.Message) error {
	opts := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	data, err := opts.Marshal(m)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
