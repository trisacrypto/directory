package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg"
	profiles "github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store/wire"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var profile *profiles.Profile

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "gdsutil"
	app.Version = pkg.Version()
	app.Usage = "utilities for operating the GDS service and database"
	app.Before = loadProfile
	app.Commands = []*cli.Command{
		{
			Name:     "ldb:populate",
			Usage:    "populate empty audit logs based on current VASP data",
			Category: "leveldb",
			Action:   populate,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
			},
		},
		{
			Name:     "ldb:keys",
			Usage:    "list the keys currently in the leveldb store",
			Category: "leveldb",
			Action:   ldbKeys,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "b64encode",
					Aliases: []string{"b"},
					Usage:   "base64 encode keys (otherwise they will be utf-8 decoded)",
				},
				&cli.StringFlag{
					Name:    "prefix",
					Aliases: []string{"p"},
					Usage:   "specify a prefix to filter keys on",
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
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "b64encode",
					Aliases: []string{"b"},
					Usage:   "specify the keys as base64 encoded values which must be decoded",
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "write the fetched key to directory if specified, otherwise printed",
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
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "b64encode",
					Aliases: []string{"b"},
					Usage:   "specify the key and value as base64 encoded strings which must be decoded",
				},
				&cli.StringFlag{
					Name:    "format",
					Aliases: []string{"f"},
					Usage:   "format of the data (raw, json, pb)",
					Value:   "json",
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
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "b64encode",
					Aliases: []string{"b"},
					Usage:   "specify the keys as base64 encoded values which must be decoded",
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
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write CSV data out to",
					Value:   "directory.csv",
				},
			},
		},
		{
			Name:      "profile",
			Aliases:   []string{"config"},
			Usage:     "view and manage profiles to configure gdsutil with",
			UsageText: "gdsutil profile [name]\n   gdsutil profile --activate [name]\n   gdsutil profile --list\n   gdsutil profile --path\n   gdsutil profile --install",
			Action:    manageProfiles,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "l",
					Aliases: []string{"list"},
					Usage:   "list the available profiles and exit",
				},
				&cli.BoolFlag{
					Name:    "p",
					Aliases: []string{"path"},
					Usage:   "show the path to the configuration and exit",
				},
				&cli.BoolFlag{
					Name:    "i",
					Aliases: []string{"install"},
					Usage:   "install the default profiles and exit",
				},
				&cli.StringFlag{
					Name:    "a",
					Aliases: []string{"activate"},
					Usage:   "activate the profile with the specified name",
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
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect the directory replica client",
					EnvVars: []string{"TRISA_DIRECTORY_REPLICA_URL"},
				},
				// TODO allow the user to add multiple peers at a time?
				&cli.Uint64Flag{
					Name:    "pid",
					Aliases: []string{"p"},
					Usage:   "specify the pid for the peer to add",
				},
				&cli.BoolFlag{
					Name:    "S",
					Aliases: []string{"no-secure"},
					Usage:   "do not connect via TLS (e.g. for development)",
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
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect the directory replica client",
					EnvVars: []string{"TRISA_DIRECTORY_REPLICA_URL"},
				},
				// TODO allow the user to rm multiple peers at a time?
				&cli.Uint64Flag{
					Name:    "pid",
					Aliases: []string{"p"},
					Usage:   "specify the pid for the peer to tombstone",
				},
				&cli.BoolFlag{
					Name:    "S",
					Aliases: []string{"no-secure"},
					Usage:   "do not connect via TLS (e.g. for development)",
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
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect the directory replica client",
					EnvVars: []string{"TRISA_DIRECTORY_REPLICA_URL"},
				},
				// TODO: have we standardized on how to reference regions?
				&cli.StringSliceFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "specify a region for peers to be returned",
				},
				&cli.BoolFlag{
					Name:    "status",
					Aliases: []string{"s"},
					Usage:   "specify for status-only, will not return peer details",
				},
				&cli.BoolFlag{
					Name:    "S",
					Aliases: []string{"no-secure"},
					Usage:   "do not connect via TLS (e.g. for development)",
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
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "partial",
					Aliases: []string{"p"},
					Usage:   "ignore any objects not specified in request",
				},
				&cli.StringSliceFlag{
					Name:    "namespaces",
					Aliases: []string{"n"},
					Usage:   "specify the namespaces to replicate (if empty, all are replicated)",
				},
				&cli.StringSliceFlag{
					Name:    "objects",
					Aliases: []string{"o"},
					Usage:   "specify the object keys to replicate (otherwise all objects from namespaces will be used)",
				},
				&cli.BoolFlag{
					Name:    "dryrun",
					Aliases: []string{"D"},
					Usage:   "show changes that would occur, does not modify database",
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
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.StringFlag{
					Name:    "addr",
					Aliases: []string{"a"},
					Usage:   "bind addr of the local replica (for name processing)",
					EnvVars: []string{"GDS_REPLICA_BIND_ADDR"},
				},
				&cli.Uint64Flag{
					Name:    "pid",
					Aliases: []string{"p"},
					Usage:   "process id of the local replica",
					EnvVars: []string{"GDS_REPLICA_PID"},
				},
				&cli.StringFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "geographic region of the local replica",
					EnvVars: []string{"GDS_REPLICA_REGION"},
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "human readable name of the local replica",
					EnvVars: []string{"GDS_REPLICA_NAME"},
				},
				&cli.BoolFlag{
					Name:    "dryrun",
					Aliases: []string{"D"},
					Usage:   "show changes that would occur, does not modify database",
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
				&cli.StringFlag{
					Name:    "key",
					Aliases: []string{"k"},
					Usage:   "secret key to decrypt the cipher text",
					EnvVars: []string{"GDS_SECRET_KEY"},
				},
			},
		},
		{
			Name:     "register:export",
			Usage:    "export a registration form for the GDS UI (e.g. to submit from TestNet to prod)",
			Category: "admin",
			Action:   registerExport,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.StringFlag{
					Name:    "id",
					Aliases: []string{"i"},
					Usage:   "VASP ID to lookup registration",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "VASP Name (common name) to lookup registration",
				},
				&cli.StringFlag{
					Name:    "admin-endpoint",
					Aliases: []string{"e"},
					Usage:   "endpoint to export registration for",
				},
				&cli.StringFlag{
					Name:    "common-name",
					Aliases: []string{"c"},
					Usage:   "common name to export registration for",
				},
				&cli.StringFlag{
					Name:    "outpath",
					Aliases: []string{"o"},
					Usage:   "path to write out JSON form to",
				},
			},
		},
		{
			Name:     "register:repair",
			Usage:    "attempt to repair a certificate request interactively",
			Category: "admin",
			Action:   registerRepair,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.StringFlag{
					Name:    "id",
					Aliases: []string{"i"},
					Usage:   "VASP ID to lookup registration",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "VASP Name (common name) to lookup registration",
				},
			},
		},
		{
			Name:     "register:reissue",
			Usage:    "create a new certificate request for the VASP",
			Category: "admin",
			Action:   registerReissue,
			Before:   openLevelDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to trisa directory storage",
					EnvVars: []string{"GDS_DATABASE_URL"},
				},
				&cli.StringFlag{
					Name:    "id",
					Aliases: []string{"i"},
					Usage:   "VASP ID to lookup registration",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "VASP Name (common name) to lookup registration",
				},
				&cli.StringFlag{
					Name:    "reason",
					Aliases: []string{"r"},
					Usage:   "reason for reissuing the certificates",
				},
				&cli.StringFlag{
					Name:    "email",
					Aliases: []string{"e"},
					Usage:   "email of user reissuing certs for audit log",
				},
			},
		},
		{
			Name:     "admin:tokenkey",
			Usage:    "generate an RSA token key pair and ksuid for JWT token signing",
			Category: "admin",
			Action:   generateTokenKey,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write keys out to (optional, will be saved as ksuid.pem by default)",
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "number of bits for the generated keys",
					Value:   4096,
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
		return cli.Exit(err, 1)
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
			return cli.Exit("specify an existing, writeable directory to output files to", 1)
		}
		if !info.IsDir() {
			return cli.Exit("specify a directory to write files out to", 1)
		}
	}

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args().Slice() {
		var key []byte
		if key, err = wire.DecodeKey(keys, b64decode); err != nil {
			return cli.Exit(fmt.Errorf("could not decode key: %s", err), 1)
		}

		var data []byte
		if data, err = ldb.Get(key, nil); err != nil {
			return cli.Exit(err, 1)
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
				return cli.Exit(err, 1)
			}
		case wire.NamespaceIndices:
			if jsonValue, err = wire.UnmarshalIndex(data); err != nil {
				return cli.Exit(err, 1)
			}
		case wire.NamespaceSequence:
			if jsonValue, err = wire.UnmarshalSequence(data); err != nil {
				return cli.Exit(err, 1)
			}
		default:
			fmt.Fprintf(os.Stderr, "warning: cannot unmarshal unknown namespace %q, printing raw data\n", prefix)
		}

		// Marshal JSON representation for pretty-printing
		var outdata []byte
		switch {
		case jsonValue != nil:
			if outdata, err = json.MarshalIndent(jsonValue, "", "  "); err != nil {
				return cli.Exit(err, 1)
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
				return cli.Exit(err, 1)
			}
		default:
			outdata = data
		}

		if out != "" {
			path := filepath.Join(out, string(key)+".json")
			if err = ioutil.WriteFile(path, outdata, 0644); err != nil {
				return cli.Exit(err, 1)
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
		return cli.Exit("specify path, key and path, or key and value as arguments", 1)
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
			return cli.Exit(fmt.Errorf("mismatch file extension %q and data format %q: specify --format", ext, format), 1)
		}

		key = []byte(strings.TrimSuffix(name, ext))
		if data, err = ioutil.ReadFile(path); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		key = []byte(args.Get(0))

		// Determine if second argument is a path
		varg := args.Get(1)
		if isFile(varg) {
			ext := filepath.Ext(varg)
			if strings.TrimLeft(ext, ".") != format {
				return cli.Exit(fmt.Errorf("mismatch file extension %q and data format %q: specify --format", ext, format), 1)
			}
			if data, err = ioutil.ReadFile(varg); err != nil {
				return cli.Exit(err, 1)
			}
		} else {
			data = []byte(varg)
		}

	}

	// Perform base64 decoding if necessary
	b64decode := c.Bool("b64decode")
	if b64decode {
		if key, err = base64.RawStdEncoding.DecodeString(string(key)); err != nil {
			return cli.Exit(err, 1)
		}
		if data, err = base64.RawStdEncoding.DecodeString(string(data)); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Quick spot check
	if len(data) == 0 || len(key) == 0 {
		return cli.Exit("no key or value found", 1)
	}

	// Unmarshal the thing from data then
	// Marshal the database representation
	if format != "raw" && format != "bytes" {
		// Prefix is required on key to determine how to unmarshal the data
		prefix := strings.Split(string(key), "::")[0]
		switch format {
		case "json":
			if value, err = wire.RemarshalJSON(prefix, data); err != nil {
				return cli.Exit(err, 1)
			}
		case "pb", "proto", "protobuf":
			// Check if the protocol buffers can be unmarshaled; if so, the data is good to go
			if _, err = wire.UnmarshalProto(prefix, data); err != nil {
				return cli.Exit(err, 1)
			}
			value = data
		default:
			return cli.Exit("unknown format: specify raw, bytes, json, or proto", 1)
		}

	} else {
		// Raw or bytes data is just the data
		value = data
	}

	// Final spot check
	if len(value) == 0 {
		return cli.Exit("no value marshaled", 1)
	}

	// Put the key/value to the database
	if err = ldb.Put(key, value, nil); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func ldbDelete(c *cli.Context) (err error) {
	defer ldb.Close()
	if c.NArg() == 0 {
		return cli.Exit("specify at least one key to delete", 1)
	}

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args().Slice() {
		var key []byte
		if key, err = wire.DecodeKey(keys, b64decode); err != nil {
			return cli.Exit(err, 1)
		}

		if err = ldb.Delete(key, nil); err != nil {
			return cli.Exit(err, 1)
		}
	}

	return nil
}

func populate(c *cli.Context) (err error) {
	defer ldb.Close()

	populated := 0

	iter := ldb.NewIterator(util.BytesPrefix([]byte(wire.NamespaceVASPs)), nil)
	for iter.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(iter.Value(), vasp); err != nil {
			iter.Release()
			fmt.Printf("populated %d VASP logs\n", populated)
			return cli.Exit(err, 1)
		}

		var log []*models.AuditLogEntry
		if log, err = models.GetAuditLog(vasp); err != nil {
			iter.Release()
			fmt.Printf("populated %d VASP logs\n", populated)
			return cli.Exit(err, 1)
		}

		if len(log) > 0 {
			continue
		}

		// Fill in initial state
		var entry *models.AuditLogEntry

		switch {
		case vasp.VerificationStatus == pb.VerificationState_NO_VERIFICATION:
			// Initial state, I don't think this can happen in the current code
			// because new VASPs get automatically set to SUBMITTED but I'm adding
			// the case here for completeness.
			entry = &models.AuditLogEntry{
				Timestamp:     vasp.LastUpdated,
				PreviousState: pb.VerificationState_NO_VERIFICATION,
				CurrentState:  pb.VerificationState_NO_VERIFICATION,
				Description:   "first listed",
				Source:        "automated",
			}
			if err = models.AppendAuditLog(vasp, entry); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}
		case vasp.VerificationStatus == pb.VerificationState_SUBMITTED:
			// VASP was only submitted
			entry = &models.AuditLogEntry{
				Timestamp:     vasp.FirstListed,
				PreviousState: pb.VerificationState_NO_VERIFICATION,
				CurrentState:  pb.VerificationState_SUBMITTED,
				Description:   "register request received",
				Source:        "automated",
			}
			if err = models.AppendAuditLog(vasp, entry); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}
		case vasp.VerifiedOn == "":
			// VASP was never verified

			// Add initial entry
			entry = &models.AuditLogEntry{
				Timestamp:     vasp.FirstListed,
				PreviousState: pb.VerificationState_NO_VERIFICATION,
				CurrentState:  pb.VerificationState_SUBMITTED,
				Description:   "register request received",
				Source:        "automated",
			}
			if err = models.AppendAuditLog(vasp, entry); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}

			// Add entry for the current state
			entry = &models.AuditLogEntry{
				Timestamp:     vasp.LastUpdated,
				PreviousState: pb.VerificationState_SUBMITTED,
				CurrentState:  vasp.VerificationStatus,
				Description:   "reconstructed VASP state change",
				Source:        "automated",
			}
			if err = models.AppendAuditLog(vasp, entry); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}
		case vasp.VerifiedOn != "":
			// VASP was verified at some point

			// Add intital entry
			entry = &models.AuditLogEntry{
				Timestamp:     vasp.FirstListed,
				PreviousState: pb.VerificationState_NO_VERIFICATION,
				CurrentState:  pb.VerificationState_SUBMITTED,
				Description:   "register request received",
				Source:        "automated",
			}
			if err = models.AppendAuditLog(vasp, entry); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}

			// Add verified entry
			entry = &models.AuditLogEntry{
				Timestamp:     vasp.VerifiedOn,
				PreviousState: pb.VerificationState_PENDING_REVIEW,
				CurrentState:  pb.VerificationState_REVIEWED,
				Description:   "registration request received",
				Source:        "automated",
			}
			if err = models.AppendAuditLog(vasp, entry); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}

			if vasp.VerificationStatus != pb.VerificationState_REVIEWED {
				// VASP changed to another state
				entry = &models.AuditLogEntry{
					Timestamp:     vasp.LastUpdated,
					PreviousState: pb.VerificationState_REVIEWED,
					CurrentState:  vasp.VerificationStatus,
					Description:   "reconstructed VASP state change",
					Source:        "automated",
				}
				if err = models.AppendAuditLog(vasp, entry); err != nil {
					iter.Release()
					fmt.Printf("populated %d VASP logs\n", populated)
					return cli.Exit(err, 1)
				}
			}
		default:
			fmt.Printf("unknown state for VASP %s\n", vasp.Id)
		}

		if entry != nil {
			key := makeKey([]byte("vasps::"), vasp.Id)
			var value []byte
			if value, err = proto.Marshal(vasp); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}

			// Put the key/value to the database
			if err = ldb.Put(key, value, nil); err != nil {
				iter.Release()
				fmt.Printf("populated %d VASP logs\n", populated)
				return cli.Exit(err, 1)
			}
			populated++
		}
	}

	fmt.Printf("successfully populated %d VASP logs\n", populated)

	return nil
}

func makeKey(prefix []byte, id string) (key []byte) {
	buf := []byte(id)
	key = make([]byte, 0, len(prefix)+len(buf))
	key = append(key, prefix...)
	key = append(key, buf...)
	return key
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
			return cli.Exit(err, 1)
		}

		record := make(map[string]string)
		record["common_name"] = vasp.CommonName
		record["name"], _ = vasp.Name()
		record["key"] = string(iter.Key())
		record["registered_directory"] = vasp.RegisteredDirectory
		record["vasp_status"] = vasp.VerificationStatus.String()
		record["verified_on"] = vasp.VerifiedOn
		data[vasp.Id] = record
	}

	if err = iter.Error(); err != nil {
		iter.Release()
		return cli.Exit(err, 1)
	}
	iter.Release()

	// Iterate over certreqs
	iter = ldb.NewIterator(util.BytesPrefix([]byte(wire.NamespaceCertReqs)), nil)
	for iter.Next() {
		cr := new(models.CertificateRequest)
		if err = proto.Unmarshal(iter.Value(), cr); err != nil {
			iter.Release()
			return cli.Exit(err, 1)
		}

		record, ok := data[cr.Vasp]
		if !ok {
			fmt.Println("no VASP for certificate request", string(iter.Key()))
			continue
		}
		record["certreq"] = cr.Id
		record["certreq_key"] = string(iter.Key())
		record["certreq_status"] = cr.Status.String()
	}

	if err = iter.Error(); err != nil {
		iter.Release()
		return cli.Exit(err, 1)
	}
	iter.Release()

	// Write out a CSV file of the VASP list
	var f *os.File
	if f, err = os.OpenFile(c.String("out"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return cli.Exit(err, 1)
	}
	w := csv.NewWriter(f)
	w.Write([]string{"id", "name", "common_name", "registered_directory", "verified_on", "verification_status", "certreq", "certreq_status"})
	for id, record := range data {
		row := []string{id, record["name"], record["common_name"], record["registered_directory"], record["verified_on"], record["vasp_status"], record["certreq"], record["certreq_status"]}
		w.Write(row)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("%d records written to %s\n", len(data), c.String("out"))
	return nil
}

//===========================================================================
// LevelDB Helper Functions
//===========================================================================

func openLevelDB(c *cli.Context) (err error) {
	if ldb, err = profile.OpenLevelDB(); err != nil {
		return cli.Exit(err, 1)
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
// Profile Actions
//===========================================================================

func manageProfiles(c *cli.Context) (err error) {
	// Handle list and then exit
	if c.Bool("list") {
		var p *profiles.Profiles
		if p, err = profiles.Load(); err != nil {
			return cli.Exit(err, 1)
		}

		if len(p.Profiles) == 0 {
			fmt.Println("no available profiles")
			return nil
		}

		fmt.Println("available profiles\n------------------")
		for name := range p.Profiles {
			if name == p.Active {
				fmt.Printf("- *%s\n", name)
			} else {
				fmt.Printf("-  %s\n", name)
			}

		}

		return nil
	}

	// Handle path and then exit
	if c.Bool("path") {
		var path string
		if path, err = profiles.ProfilesPath(); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Println(path)
		return nil
	}

	// Handle install and then exit
	if c.Bool("install") {
		if err = profiles.Install(); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	// Handle activate and then exit
	if name := c.String("activate"); name != "" {
		var p *profiles.Profiles
		if p, err = profiles.Load(); err != nil {
			return cli.Exit(err, 1)
		}

		if err = p.SetActive(name); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Printf("profile %q is now active\n", name)
		return nil
	}

	// Handle show named or active profile
	if c.Args().Len() > 1 {
		return cli.Exit("specify only a single profile to print", 1)
	}
	var p *profiles.Profiles
	if p, err = profiles.Load(); err != nil {
		return cli.Exit(err, 1)
	}

	if profile, err = p.GetActive(c.Args().Get(0)); err != nil {
		return cli.Exit(err, 1)
	}

	var data []byte
	if data, err = yaml.Marshal(p.Profiles[p.Active]); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(data))
	return nil
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
	ctx, cancel := profile.Context()
	defer cancel()

	// call client.AddPeer with the pid
	var out *peers.PeersStatus
	if out, err = replicaClient.AddPeers(ctx, peer); err != nil {
		return cli.Exit(err, 1)
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
	ctx, cancel := profile.Context()
	defer cancel()

	// call client.RmPeer with the pid
	var out *peers.PeersStatus
	if out, err = replicaClient.RmPeers(ctx, peer); err != nil {
		return cli.Exit(err, 1)
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
	ctx, cancel := profile.Context()
	defer cancel()

	// call client.GetPeers with filter
	var out *peers.PeersList
	if out, err = replicaClient.GetPeers(ctx, filter); err != nil {
		return cli.Exit(err, 1)
	}

	// print the peers
	printJSON(out)
	return nil
}

func initReplicaClient(c *cli.Context) (err error) {
	if replicaClient, err = profile.Replica.Connect(); err != nil {
		return cli.Exit(err, 1)
	}
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
		return cli.Exit("must specify ciphertext and hmac arguments", 1)
	}

	var secret string
	if secret = c.String("key"); secret == "" {
		return cli.Exit("cipher key required", 1)
	}

	var ciphertext, signature []byte
	if ciphertext, err = base64.RawStdEncoding.DecodeString(c.Args().Get(0)); err != nil {
		return cli.Exit(fmt.Errorf("could not decode ciphertext: %s", err), 1)
	}
	if signature, err = base64.RawStdEncoding.DecodeString(c.Args().Get(1)); err != nil {
		return cli.Exit(fmt.Errorf("could not decode signature: %s", err), 1)
	}

	if len(ciphertext) == 0 {
		return cli.Exit("empty cipher text", 1)
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
		return cli.Exit(err, 1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return cli.Exit(err, 1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return cli.Exit(err, 1)
	}

	plainbytes, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return cli.Exit(err, 1)
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
// Admin Functions
//===========================================================================

func registerExport(c *cli.Context) (err error) {
	defer ldb.Close()

	vaspID := c.String("id")
	name := c.String("name")

	// Lookup VASP in database by ID or by name
	var vasp *pb.VASP
	switch {
	case vaspID != "":
		if vasp, err = getVASPByID(vaspID); err != nil {
			return cli.Exit(err, 1)
		}
	case name != "":
		if vasp, err = getVASPByCommonName(name); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		return cli.Exit("specify either ID or common name for lookup", 1)
	}

	// Remove sensitive data from contacts
	for _, contact := range []*pb.Contact{vasp.Contacts.Technical, vasp.Contacts.Administrative, vasp.Contacts.Legal, vasp.Contacts.Billing} {
		if contact != nil {
			contact.Extra = nil
		}
	}

	pbForm := &api.RegisterRequest{
		Entity:           vasp.Entity,
		Contacts:         vasp.Contacts,
		TrisaEndpoint:    c.String("directory-endpoint"),
		CommonName:       c.String("common-name"),
		Website:          vasp.Website,
		BusinessCategory: vasp.BusinessCategory,
		VaspCategories:   vasp.VaspCategories,
		EstablishedOn:    vasp.EstablishedOn,
		Trixo:            vasp.Trixo,
	}

	// Intermediate marshal then unmarshal ensures that all fields are exported even
	// if they are empty (so the front-end UI doesn't break on upload).
	jsonpb := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}

	data, err := jsonpb.Marshal(pbForm)
	if err != nil {
		return cli.Exit(err, 1)
	}

	registrationForm := make(map[string]interface{})
	if err = json.Unmarshal(data, &registrationForm); err != nil {
		return cli.Exit(err, 1)
	}

	form := map[string]interface{}{
		"version":          "v1beta1",
		"registrationForm": registrationForm,
	}

	var w io.Writer
	if path := c.String("outpath"); path != "" {
		var f *os.File
		if f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return cli.Exit(err, 1)
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(form); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func registerRepair(c *cli.Context) (err error) {
	defer ldb.Close()

	vaspID := c.String("id")
	name := c.String("name")

	// Lookup VASP in database by ID or by name
	var vasp *pb.VASP
	switch {
	case vaspID != "":
		if vasp, err = getVASPByID(vaspID); err != nil {
			return cli.Exit(err, 1)
		}
	case name != "":
		if vasp, err = getVASPByCommonName(name); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		return cli.Exit("specify either ID or common name for lookup", 1)
	}

	// Find the CertificateRequest for the VASP
	var certreq *models.CertificateRequest
	if certreq, err = findCertificateRequest(vasp.Id); err != nil {
		return cli.Exit(err, 1)
	}

	if certreq == nil {
		fmt.Println("VASP has no certificate request: creating new request with new PKCS12 password")

		var conf config.Config
		if conf, err = config.New(); err != nil {
			return cli.Exit(err, 1)
		}

		// Connect to secret manager
		var sm *secrets.SecretManager
		if sm, err = secrets.New(conf.Secrets); err != nil {
			return cli.Exit(err, 1)
		}

		// Create PKCS12 password along with certificate request.
		password := secrets.CreateToken(16)
		certreq = &models.CertificateRequest{
			Id:         uuid.New().String(),
			Vasp:       vasp.Id,
			CommonName: vasp.CommonName,
			Status:     models.CertificateRequestState_INITIALIZED,
			Created:    time.Now().Format(time.RFC3339),
		}

		// Make a new secret of type "password"
		secretType := "password"
		if err = sm.With(certreq.Id).CreateSecret(context.TODO(), secretType); err != nil {
			return cli.Exit(err, 1)
		}
		if err = sm.With(certreq.Id).AddSecretVersion(context.TODO(), secretType, []byte(password)); err != nil {
			return cli.Exit(err, 1)
		}

		var data []byte
		certreq.Modified = time.Now().Format(time.RFC3339)
		key := []byte(wire.NamespaceCertReqs + "::" + certreq.Id)
		if data, err = proto.Marshal(certreq); err != nil {
			return cli.Exit(err, 1)
		}

		if err = ldb.Put(key, data, nil); err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Printf("created new certificate request: %s\n", key)
		fmt.Printf("pkcs12 password: %s\n", password)
	}

	return nil
}

func registerReissue(c *cli.Context) (err error) {
	defer ldb.Close()

	vaspID := c.String("id")
	name := c.String("name")
	reason := c.String("reason")
	email := c.String("email")

	// Make sure there is a reason
	if reason == "" || email == "" {
		return cli.Exit("supply a reason and email of user to reissue the certs", 1)
	}

	// Lookup VASP in database by ID or by name
	var vasp *pb.VASP
	switch {
	case vaspID != "":
		if vasp, err = getVASPByID(vaspID); err != nil {
			return cli.Exit(err, 1)
		}
	case name != "":
		if vasp, err = getVASPByCommonName(name); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		return cli.Exit("specify either ID or common name for lookup", 1)
	}

	// Find the current CertificateRequest for the VASP
	var certreq *models.CertificateRequest
	if certreq, err = findCertificateRequest(vasp.Id); err != nil {
		return cli.Exit(err, 1)
	}

	// Update the current CertificateRequest if it exists
	if certreq != nil {
		// Check the current certreq status; if it hasn't already been downloaded, then cancel it.
		if certreq.Status < models.CertificateRequestState_COMPLETED {
			fmt.Printf("canceling certificate request %s and setting state %s from %s\n", certreq.Id, models.CertificateRequestState_CR_ERRORED, certreq.Status)
			if err = models.UpdateCertificateRequestStatus(certreq, models.CertificateRequestState_CR_ERRORED, reason, email); err != nil {
				return cli.Exit(err, 1)
			}
			certreq.RejectReason = reason
			certreq.Modified = time.Now().Format(time.RFC3339)

			var data []byte
			key := []byte(wire.NamespaceCertReqs + "::" + certreq.Id)
			if data, err = proto.Marshal(certreq); err != nil {
				return cli.Exit(err, 1)
			}

			if err = ldb.Put(key, data, nil); err != nil {
				return cli.Exit(err, 1)
			}
		} else {
			fmt.Printf("certificate request %s is in state %s - making no changes\n", certreq.Id, certreq.Status)
		}
	}

	// Connect to the SecretManager to create a new PKCS12 Password
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	// Connect to secret manager
	var sm *secrets.SecretManager
	if sm, err = secrets.New(conf.Secrets); err != nil {
		return cli.Exit(err, 1)
	}

	// Create a new certificate request for the VASP along with new PKCS12 password
	password := secrets.CreateToken(16)
	certreq = &models.CertificateRequest{
		Id:         uuid.New().String(),
		Vasp:       vasp.Id,
		CommonName: vasp.CommonName,
		Created:    time.Now().Format(time.RFC3339),
	}

	if err = models.UpdateCertificateRequestStatus(certreq, models.CertificateRequestState_READY_TO_SUBMIT, "reissue certificates", email); err != nil {
		return cli.Exit(err, 1)
	}

	// Make a new secret of type "password"
	secretType := "password"
	if err = sm.With(certreq.Id).CreateSecret(context.TODO(), secretType); err != nil {
		return cli.Exit(err, 1)
	}
	if err = sm.With(certreq.Id).AddSecretVersion(context.TODO(), secretType, []byte(password)); err != nil {
		return cli.Exit(err, 1)
	}

	var data []byte
	certreq.Modified = time.Now().Format(time.RFC3339)
	key := []byte(wire.NamespaceCertReqs + "::" + certreq.Id)
	if data, err = proto.Marshal(certreq); err != nil {
		return cli.Exit(err, 1)
	}

	if err = ldb.Put(key, data, nil); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("created new certificate request: %s\n", key)
	fmt.Printf("pkcs12 password: %s\n", password)
	return nil
}

func getVASPByID(id string) (vasp *pb.VASP, err error) {
	var value []byte
	key := []byte(fmt.Sprintf("vasps::%s", id))

	if value, err = ldb.Get(key, nil); err != nil {
		return nil, err
	}

	vasp = new(pb.VASP)
	if err = proto.Unmarshal(value, vasp); err != nil {
		return nil, err
	}

	return vasp, nil
}

func getVASPByCommonName(name string) (_ *pb.VASP, err error) {
	var names []byte
	if names, err = ldb.Get([]byte("index::names"), nil); err != nil {
		return nil, err
	}

	var index map[string]interface{}
	if index, err = wire.UnmarshalIndex(names); err != nil {
		return nil, err
	}

	if id, ok := index[name]; ok {
		return getVASPByID(id.(string))
	}

	return nil, fmt.Errorf("couldn't find VASP with common name %q", name)
}

func findCertificateRequest(vaspID string) (cr *models.CertificateRequest, err error) {
	iter := ldb.NewIterator(util.BytesPrefix([]byte(wire.NamespaceCertReqs)), nil)
	defer iter.Release()
	for iter.Next() {
		cr = new(models.CertificateRequest)
		if err = proto.Unmarshal(iter.Value(), cr); err != nil {
			return nil, err
		}

		if cr.Vasp == vaspID {
			return cr, nil
		}
	}

	if err = iter.Error(); err != nil {
		return nil, err
	}

	// Couldn't find the certificate request, but don't return an error
	return nil, nil
}

func generateTokenKey(c *cli.Context) (err error) {
	// Create ksuid and determine outpath
	var keyid ksuid.KSUID
	if keyid, err = ksuid.NewRandom(); err != nil {
		return cli.Exit(err, 1)
	}

	var out string
	if out = c.String("out"); out == "" {
		out = fmt.Sprintf("%s.pem", keyid)
	}

	// Generate RSA keys using crypto random
	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, c.Int("size")); err != nil {
		return cli.Exit(err, 1)
	}

	// Open file to PEM encode keys to
	var f *os.File
	if f, err = os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return cli.Exit(err, 1)
	}

	if err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("RSA key id: %s -- saved with PEM encoding to %s\n", keyid, out)
	return nil
}

//===========================================================================
// Helper Functions
//===========================================================================

// loadProfile runs before every command so it cannot return an error; if it cannot
// load the profile, it will attempt to create a default profile unless a named profile
// was given.
func loadProfile(c *cli.Context) (err error) {
	if profile, err = profiles.LoadActive(c); err != nil {
		if name := c.String("profile"); name != "" {
			return cli.Exit(err, 1)
		}
		profile = profiles.New()
		if err = profile.Update(c); err != nil {
			return cli.Exit(err, 1)
		}
	}
	return nil
}

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
		return cli.Exit(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
