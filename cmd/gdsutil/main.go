package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/store/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
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
					Name:  "b, b64encode",
					Usage: "base64 encode keys (otherwise they will be utf-8 decoded)",
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
			Name:     "ldb:list",
			Usage:    "print a summary of the current contents of the database",
			Category: "leveldb",
			Action:   ldbList,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
			},
		},
		{
			Name:     "peers:add",
			Usage:    "add peers to the network by pid",
			Category: "replica",
			Before:   initReplicaClient,
			Action:   addPeers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "u, replica-endpoint",
					Usage:  "the url to connect the directory replica client",
					Value:  "replica.vaspdirectory.net:443",
					EnvVar: "TRISA_DIRECTORY_REPLICA_URL",
				},
				// TODO allow the user to add multiple peers at a time?
				cli.Uint64Flag{
					Name:  "p, pid",
					Usage: "specify the pid for the peer to add",
				},
			},
		},
		{
			Name:     "peers:delete",
			Usage:    "remove a peer from the network by pid",
			Category: "replica",
			Before:   initReplicaClient,
			Action:   delPeers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "u, replica-endpoint",
					Usage:  "the url to connect the directory replica client",
					Value:  "replica.vaspdirectory.net:443",
					EnvVar: "TRISA_DIRECTORY_REPLICA_URL",
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
			Before:   initReplicaClient,
			Action:   listPeers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "u, replica-endpoint",
					Usage:  "the url to connect the directory replica client",
					Value:  "replica.vaspdirectory.net:443",
					EnvVar: "TRISA_DIRECTORY_REPLICA_URL",
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

var ldb *leveldb.DB

func ldbKeys(c *cli.Context) (err error) {
	defer ldb.Close()

	var prefix *util.Range
	if prefixs := c.String("prefix"); prefixs != "" {
		prefix = util.BytesPrefix([]byte(prefixs))
	}

	iter := ldb.NewIterator(prefix, nil)
	defer iter.Release()

	b64encode := c.Bool("b64encode")
	for iter.Next() {
		fmt.Printf("- %s\n", wire.EncodeKey(iter.Key(), b64encode))
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
			return cli.NewExitError("specify an existing, writeable directory to output files to", 1)
		}
		if !info.IsDir() {
			return cli.NewExitError("specify a directory to write files out to", 1)
		}
	}

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args() {
		var key []byte
		if key, err = wire.DecodeKey(keys, b64decode); err != nil {
			return cli.NewExitError(fmt.Errorf("could not decode key: %s", err), 1)
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

		prefix := strings.Split(keys, "::")[0]
		switch prefix {
		case wire.NamespaceVASPs, wire.NamespaceCertReqs, wire.NamespaceReplicas:
			if pbValue, err = wire.UnmarshalProto(prefix, data); err != nil {
				return cli.NewExitError(err, 1)
			}
		case wire.NamespaceIndices:
			if jsonValue, err = wire.UnmarshalIndex(data); err != nil {
				return cli.NewExitError(err, 1)
			}
		case wire.NamespaceSequence:
			if jsonValue, err = wire.UnmarshalSequence(data); err != nil {
				return cli.NewExitError(err, 1)
			}
		default:
			fmt.Fprintf(os.Stderr, "warning: cannot unmarshal unknown namespace %q, printing raw data\n", prefix)
		}

		// Marshal JSON representation for pretty-printing
		var outdata []byte
		switch {
		case jsonValue != nil:
			if outdata, err = json.MarshalIndent(jsonValue, "", "  "); err != nil {
				return cli.NewExitError(err, 1)
			}
		case pbValue != nil:
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
		default:
			outdata = data
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
		// Prefix is required on key to determine how to unmarshal the data
		prefix := strings.Split(string(key), "::")[0]
		switch format {
		case "json":
			if value, err = wire.RemarshalJSON(prefix, data); err != nil {
				return cli.NewExitError(err, 1)
			}
		case "pb", "proto", "protobuf":
			// Check if the protocol buffers can be unmarshaled; if so, the data is good to go
			if _, err = wire.UnmarshalProto(prefix, data); err != nil {
				return cli.NewExitError(err, 1)
			}
			value = data
		default:
			return cli.NewExitError("unknown format: specify raw, bytes, json, or proto", 1)
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
		if key, err = wire.DecodeKey(keys, b64decode); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err = ldb.Delete(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

func ldbList(c *cli.Context) (err error) {
	defer ldb.Close()

	var data = make(map[string]map[string]string)
	var iter iterator.Iterator

	// Iterate over vasps
	iter = ldb.NewIterator(util.BytesPrefix([]byte(wire.NamespaceVASPs)), nil)
	for iter.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(iter.Value(), vasp); err != nil {
			iter.Release()
			return cli.NewExitError(err, 1)
		}

		record := make(map[string]string)
		record["common_name"] = vasp.CommonName
		record["name"], _ = vasp.Name()
		record["key"] = string(iter.Key())
		record["registered_directory"] = vasp.RegisteredDirectory
		record["vasp_status"] = vasp.VerificationStatus.String()
		data[vasp.Id] = record
	}

	if err = iter.Error(); err != nil {
		iter.Release()
		return cli.NewExitError(err, 1)
	}
	iter.Release()

	// Iterate over certreqs
	iter = ldb.NewIterator(util.BytesPrefix([]byte(wire.NamespaceCertReqs)), nil)
	for iter.Next() {
		cr := new(models.CertificateRequest)
		if err = proto.Unmarshal(iter.Value(), cr); err != nil {
			iter.Release()
			return cli.NewExitError(err, 1)
		}

		record, ok := data[cr.Vasp]
		if !ok {
			fmt.Println("no VASP for certificate request", string(iter.Key()))
			continue
		}
		record["certreq"] = cr.Id
		record["certreq_key"] = string(iter.Key())
	}

	if err = iter.Error(); err != nil {
		iter.Release()
		return cli.NewExitError(err, 1)
	}
	iter.Release()

	var out []byte
	if out, err = yaml.Marshal(data); err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(string(out))
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

//===========================================================================
// Peer Management Replica Actions
//===========================================================================

var replicaClient peers.PeerManagementClient

func addPeers(c *cli.Context) (err error) {
	// Create the Peer with the specified PID
	// TODO: how to add the other values for a Peer?
	peer := &peers.Peer{
		Id: c.Uint64("pid"),
	}

	// create a new context and pass the parent context in
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// call client.AddPeer with the pid
	var out *peers.PeersStatus
	if out, err = replicaClient.AddPeers(ctx, peer); err != nil {
		return cli.NewExitError(err, 1)
	}

	// print the returned result
	printJSON(out)
	return nil
}

func delPeers(c *cli.Context) (err error) {
	peer := &peers.Peer{
		Id: c.Uint64("pid"),
	}

	// create a new context and pass the parent context in
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// call client.RmPeer with the pid
	var out *peers.PeersStatus
	if out, err = replicaClient.RmPeers(ctx, peer); err != nil {
		return cli.NewExitError(err, 1)
	}

	// print the returned result
	printJSON(out)
	return nil
}

func listPeers(c *cli.Context) (err error) {
	// determine if this is a region-specific or status only request
	filter := &peers.PeersFilter{
		Region:     c.StringSlice("region"),
		StatusOnly: c.Bool("status"),
	}

	// create a new context and pass the parent context in
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// call client.GetPeers with filter
	var out *peers.PeersList
	if out, err = replicaClient.GetPeers(ctx, filter); err != nil {
		return cli.NewExitError(err, 1)
	}

	// print the peers
	printJSON(out)
	return nil
}

func initReplicaClient(c *cli.Context) (err error) {
	// initialize a client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.String("replica-endpoint"), grpc.WithInsecure()); err != nil {
		return cli.NewExitError(err, 1)
	}
	replicaClient = peers.NewPeerManagementClient(cc)
	return nil
}

//===========================================================================
// Replica Gossip Actions
//===========================================================================

func gossip(c *cli.Context) (err error) {
	return errors.New("honu replication required")
}

func gossipMigrate(c *cli.Context) (err error) {
	return errors.New("honu object migration required")
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
