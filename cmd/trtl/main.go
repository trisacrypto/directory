package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/rotationalio/honu"
	opts "github.com/rotationalio/honu/options"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/trisacrypto/directory/pkg"
	profiles "github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	grpc "google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var profile *profiles.Profile

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "trtl"
	app.Version = pkg.Version()
	app.Usage = "a command line tool for interacting with the trtl data replication service"
	app.Before = loadProfile
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "trtl-endpoint",
			Aliases: []string{"u"},
			Usage:   "the url to connect to the trtl replication service",
			EnvVars: []string{"TRISA_DIRECTORY_REPLICA_URL", "TRTL_ENDPOINT", "TRTL_URL"},
		},
		&cli.BoolFlag{
			Name:    "no-secure",
			Aliases: []string{"S"},
			Usage:   "do not connect via TLS (e.g. for development)",
		},
		&cli.IntFlag{
			Name:    "trtl-index",
			Aliases: []string{"i"},
			Usage:   "the index (starting from 0) of the desired trtl replica",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "run the trtl database and replication service",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to start the trtl database on",
					EnvVars: []string{"TRTL_DATABASE_URL"},
				},
				&cli.StringFlag{
					Name:    "bindaddr",
					Aliases: []string{"a"},
					Usage:   "address to bind the trtl server to",
				},
				&cli.Uint64Flag{
					Name:    "pid",
					Aliases: []string{"p"},
					Usage:   "processor ID for the trtl node",
				},
				&cli.StringFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "region for the trtl node",
				},
			},
		},
		{
			Name:     "validate",
			Usage:    "validate the current trtl configuration",
			Category: "server",
			Action:   validate,
		},
		{
			Name:      "migrate",
			Usage:     "migrate a leveldb database to a trtl database",
			ArgsUsage: "src dst",
			Category:  "server",
			Action:    migrate,
		},
		{
			Name:     "migrate-certs",
			Usage:    "create certificate records in the trtl database from the existing certreqs",
			Category: "client",
			Before:   initDBClient,
			Action:   migrateCerts,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "report",
					Aliases: []string{"r"},
					Usage:   "dump a migration report to a file",
				},
				&cli.BoolFlag{
					Name:    "cleanup",
					Aliases: []string{"c"},
					Usage:   "cleanup dangling references to non-existent certificate requests",
				},
				&cli.BoolFlag{
					Name:    "timestamps",
					Aliases: []string{"t"},
					Usage:   "fill in missing verified timestamps on the VASP records",
				},
			},
		},
		{
			Name:     "restore-certs",
			Usage:    "create certificate records in the database from a sectigo audit log",
			Category: "client",
			Before:   initDBClient,
			Action:   restoreCerts,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Aliases:  []string{"p"},
					Required: true,
					Usage:    "path to the sectigo audit log",
				},
				&cli.StringFlag{
					Name:    "report",
					Aliases: []string{"r"},
					Usage:   "dump a migration report to a file",
				},
			},
		},
		{
			Name:     "status",
			Usage:    "check the status of the trtl database and replication service",
			Category: "client",
			Before:   initDBClient,
			Action:   status,
		},
		{
			Name:      "db:get",
			Usage:     "get a value from the trtl database",
			ArgsUsage: "key [key ...]",
			Category:  "client",
			Before:    initDBClient,
			Action:    dbGet,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "b64encode",
					Aliases: []string{"b"},
					Usage:   "specify the keys as base64 encoded values which must be decoded",
				},
				&cli.BoolFlag{
					Name:    "meta",
					Aliases: []string{"m"},
					Usage:   "return the metadata along with the value",
				},
				&cli.StringFlag{
					Name:    "namespace",
					Aliases: []string{"n"},
					Usage:   "specify the namespace as a string",
				},
			},
		},
		{
			Name: "db:put",
			// TODO: would we ever want to put multiple values using the CLI? If so, what arg format would we expect to use?
			Usage:    "put a single value to a key in the trtl database",
			Category: "client",
			Before:   initDBClient,
			Action:   dbPut,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "key",
					Aliases: []string{"k"},
					Usage:   "specify the key as a string",
				},
				&cli.StringFlag{
					Name:    "namespace",
					Aliases: []string{"n"},
					Usage:   "specify the namespace as a string",
				},
				&cli.StringFlag{
					Name:    "value",
					Aliases: []string{"v"},
					Usage:   "specify the value to put as a string",
				},
				&cli.BoolFlag{
					Name:    "meta",
					Aliases: []string{"m"},
					Usage:   "return the metadata along with the value",
				},
			},
		},
		{
			Name:     "db:del",
			Usage:    "delete a single key in the trtl database",
			Category: "client",
			Before:   initDBClient,
			Action:   dbDelete,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "key",
					Aliases: []string{"k"},
					Usage:   "specify the key as a string",
				},
				&cli.StringFlag{
					Name:    "namespace",
					Aliases: []string{"n"},
					Usage:   "specify the namespace as a string",
				},
				&cli.BoolFlag{
					Name:    "meta",
					Aliases: []string{"m"},
					Usage:   "return the metadata along with the value",
				},
			},
		},
		{
			Name:     "db:list",
			Usage:    "list all of the keys in the trtl database",
			Category: "client",
			Before:   initDBClient,
			Action:   dbList,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "namespace",
					Aliases: []string{"n"},
					Usage:   "specify the namespace as a string (optional)",
				},
				&cli.StringFlag{
					Name:    "prefix",
					Aliases: []string{"p"},
					Usage:   "specify a prefix of keys to list (optional)",
				},
				&cli.StringFlag{
					Name:    "seek-key",
					Aliases: []string{"s", "seek"},
					Usage:   "specify a key to seek to before iterating (optional)",
				},
				&cli.BoolFlag{
					Name:    "b64encode",
					Aliases: []string{"b"},
					Usage:   "specify the prefix/seek key as base64 encoded values which must be decoded",
				},
			},
		},
		{
			Name:     "peers:add",
			Usage:    "add peers to the network by pid",
			Category: "client",
			Before:   initPeersClient,
			Action:   addPeers,
			Flags: []cli.Flag{
				// TODO allow the user to add multiple peers at a time?
				&cli.Uint64Flag{
					Name:    "pid",
					Aliases: []string{"p"},
					Usage:   "specify the pid for the peer to add",
				},
				&cli.StringFlag{
					Name:    "addr",
					Aliases: []string{"a"},
					Usage:   "specify the addr to connect to the peer on",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "specify the name to identify the peer with",
				},
				&cli.StringFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "specify the region the peer is located in",
				},
			},
		},
		{
			Name:     "peers:delete",
			Usage:    "remove a peer from the network by pid",
			Category: "client",
			Before:   initPeersClient,
			Action:   delPeers,
			Flags: []cli.Flag{
				// TODO allow the user to rm multiple peers at a time?
				&cli.Uint64Flag{
					Name:    "pid",
					Aliases: []string{"p"},
					Usage:   "specify the pid for the peer to tombstone",
				},
			},
		},
		{
			Name:     "peers:list",
			Usage:    "get a status report of all peers in the network",
			Category: "client",
			Before:   initPeersClient,
			Action:   listPeers,
			Flags: []cli.Flag{
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
			},
		},
		{
			Name:      "gossip",
			Usage:     "initiate a gossip session with a remote replica (for debugging)",
			ArgsUsage: "remote:port",
			Category:  "client",
			Action:    gossip,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to the trtl database",
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
			Category: "client",
			Action:   gossipMigrate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Usage:   "dsn to connect to the trtl database",
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
			Name:      "profile",
			Aliases:   []string{"config", "profiles"},
			Usage:     "view and manage profiles to configure trtl with",
			UsageText: "trtl profile [name]\n   trtl profile --activate [name]\n   trtl profile --list\n   trtl profile --path\n   trtl profile --install\n   trtl profile --edit",
			Action:    manageProfiles,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list the available profiles and exit",
				},
				&cli.BoolFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "show the path to the configuration and exit",
				},
				&cli.BoolFlag{
					Name:    "install",
					Aliases: []string{"i"},
					Usage:   "install the default profiles and exit",
				},
				&cli.BoolFlag{
					Name:    "edit",
					Aliases: []string{"e"},
					Usage:   "edit the profiles YAML using $EDITOR",
				},
				&cli.StringFlag{
					Name:    "activate",
					Aliases: []string{"a"},
					Usage:   "activate the profile with the specified name",
				},
			},
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// Server Functions
//===========================================================================

// serve starts the trtl server and blocks until it is stopped.
func serve(c *cli.Context) (err error) {
	// Load the configuration from the environment
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	// Overide environment configuration from CLI flags
	if addr := c.String("bindaddr"); addr != "" {
		conf.BindAddr = addr
	}

	if pid := c.Uint64("pid"); pid > 0 {
		conf.Replica.PID = pid
	}

	if region := c.String("region"); region != "" {
		conf.Replica.Region = region
	}

	if dburl := c.String("db"); dburl != "" {
		conf.Database.URL = dburl
	}

	var server *trtl.Server
	if server, err = trtl.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = server.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

// validate checks the current trtl configuration and prints the status.
func validate(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	return printJSON(conf)
}

// migrate a leveldb database to a trtl database.
func migrate(c *cli.Context) (err error) {
	if c.NArg() != 2 {
		return cli.Exit("specify src and dst database paths", 1)
	}

	var (
		srcdb *leveldb.DB
		dstdb *honu.DB
	)

	// Open the source for reading
	if srcdb, err = leveldb.OpenFile(c.Args().Get(0), &opt.Options{ReadOnly: true}); err != nil {
		return cli.Exit(fmt.Errorf("could not open src db at %q: %s", c.Args().Get(0), err), 1)
	}

	// Open the destination for writing
	// TODO: allow user to specify replica information
	if dstdb, err = honu.Open(fmt.Sprintf("leveldb:///%s", c.Args().Get(1))); err != nil {
		return cli.Exit(fmt.Errorf("could not open dst db at %q: %s", c.Args().Get(1), err), 1)
	}

	// Loop over the source database and write into the honu database
	nKeys := 0
	iter := srcdb.NewIterator(nil, nil)
	for iter.Next() {
		// Get the key and split the namespace; this is GDS-specific logic
		parts := bytes.SplitN(iter.Key(), []byte("::"), 2)
		namespace := string(parts[0])
		key := parts[1]

		if _, err = dstdb.Put(key, iter.Value(), opts.WithNamespace(namespace)); err != nil {
			return cli.Exit(fmt.Errorf("could not put %s :: %s: %s", namespace, string(key), err), 1)
		}

		nKeys++
	}

	iter.Release()
	if err = iter.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("migrated %d keys\n", nKeys)
	return nil
}

//===========================================================================
// Initialization Functions
//===========================================================================

var dbClient pb.TrtlClient
var peersClient peers.PeerManagementClient

// initDBClient starts a trtl client with a connection to a trtl database.
func initDBClient(c *cli.Context) (err error) {
	if profile.TrtlProfiles == nil {
		return cli.Exit("no trtl profile was loaded", 1)
	}
	if len(profile.TrtlProfiles) <= c.Int("trtl-index") {
		return cli.Exit("could not find trtl profile by index", 1)
	}
	if dbClient, err = profile.TrtlProfiles[c.Int("trtl-index")].ConnectDB(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

// initPeersClient starts a trtl client with a connection to a trtl database.
func initPeersClient(c *cli.Context) (err error) {
	if profile.TrtlProfiles == nil {
		return cli.Exit("no trtl profile was loaded", 1)
	}
	if len(profile.TrtlProfiles) <= c.Int("trtl-index") {
		return cli.Exit("could not find trtl profile by index", 1)
	}
	if peersClient, err = profile.TrtlProfiles[c.Int("trtl-index")].ConnectPeers(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

//===========================================================================
// Trtl (DB) Client Functions
//===========================================================================

// migrateCerts creates certificate records in the trtl database based on the existing certreqs.
func migrateCerts(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	// Open the report file if requested
	var report *os.File
	if c.String("report") != "" {
		if report, err = os.Create(c.String("report")); err != nil {
			return cli.Exit(fmt.Errorf("could not open report file %s: %s", c.String("report"), err), 1)
		}
		defer report.Close()
	}

	// Get the list of certreqs
	req := &pb.CursorRequest{
		Namespace: wire.NamespaceVASPs,
	}
	var stream pb.Trtl_CursorClient
	if stream, err = dbClient.Cursor(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}

	var migrated, skipped int
	var results *multierror.Error

vaspLoop:
	for {
		var rep *pb.KVPair
		if rep, err = stream.Recv(); err != nil {
			if err != io.EOF {
				results = multierror.Append(results, err)
			}
			break vaspLoop
		}

		vasp := &gds.VASP{}
		if err = proto.Unmarshal(rep.Value, vasp); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		if c.Bool("timestamps") {
			if vasp.VerificationStatus == gds.VerificationState_VERIFIED && vasp.VerifiedOn == "" {
				var log []*models.AuditLogEntry
				var err error
				if log, err = models.GetAuditLog(vasp); err != nil {
					results = multierror.Append(results, err)
					continue vaspLoop
				}

			verifyLoop:
				for _, entry := range log {
					if entry.CurrentState == gds.VerificationState_VERIFIED {
						vasp.VerifiedOn = entry.Timestamp
						fmt.Printf("setting verified timestamp for %s (%s) to %s\n", vasp.Id, vasp.CommonName, vasp.VerifiedOn)
						if report != nil {
							report.WriteString(fmt.Sprintf("setting verified timestamp for %s (%s) to %s\n", vasp.Id, vasp.CommonName, vasp.VerifiedOn))
						}
						break verifyLoop
					}
				}

				if vasp.VerifiedOn == "" {
					results = multierror.Append(results, fmt.Errorf("could not find verified timestamp for %s (%s)", vasp.Id, vasp.CommonName))
				}
			}
		}

		// Get all the certificate requests for this vasp
		var ids []string
		if ids, err = models.GetCertReqIDs(vasp); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		// Get the latest completed certificate request
		var certreq *models.CertificateRequest
		var created time.Time
	certreqLoop:
		for _, id := range ids {
			var rep *pb.GetReply
			rep, err = dbClient.Get(ctx, &pb.GetRequest{
				Namespace: wire.NamespaceCertReqs,
				Key:       []byte(id),
			})
			if err != nil {
				if e, ok := grpc.FromError(err); ok {
					switch e.Code() {
					case codes.NotFound:
						results = multierror.Append(results, fmt.Errorf("could not find certreq ID %s for vasp %s (%s)", id, vasp.Id, vasp.CommonName))
						// Remove dangling certreq IDs if requested
						if c.Bool("cleanup") {
							if err = models.DeleteCertReqID(vasp, id); err != nil {
								results = multierror.Append(results, err)
							} else {
								fmt.Printf("cleaning up dangling certreq ID %s for vasp %s (%s)\n", id, vasp.Id, vasp.CommonName)
								if report != nil {
									report.WriteString(fmt.Sprintf("cleaning up dangling certreq ID %s for vasp %s (%s)\n", id, vasp.Id, vasp.CommonName))
								}
							}
						}
						continue certreqLoop
					default:
						results = multierror.Append(results, err)
					}
				} else {
					results = multierror.Append(results, err)
				}
				continue vaspLoop
			}

			// Successfully retrieved the certificate request
			current := &models.CertificateRequest{}
			if err = proto.Unmarshal(rep.Value, current); err != nil {
				results = multierror.Append(results, err)
				continue vaspLoop
			}

			// Check if this is the latest completed certificate request
			if current.Status == models.CertificateRequestState_COMPLETED {
				var currentCreated time.Time
				var err error
				if currentCreated, err = time.Parse(time.RFC3339, current.Created); err != nil {
					results = multierror.Append(results, err)
					continue vaspLoop
				}

				if certreq == nil || currentCreated.After(created) {
					certreq = current
					created = currentCreated
				}
			}
		}

		if certreq == nil {
			fmt.Printf("skipping vasp %s (%s) with no certificates\n", vasp.Id, vasp.CommonName)
			if report != nil {
				report.WriteString(fmt.Sprintf("skipping vasp %s (%s) with no certificates\n", vasp.Id, vasp.CommonName))
			}
			skipped++
			continue vaspLoop
		}

		// Create the certificate record
		var cert *models.Certificate
		if cert, err = models.NewCertificate(vasp, certreq, vasp.IdentityCertificate); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		// Validate that the serial number is correct
		expected := fmt.Sprintf("%X", vasp.IdentityCertificate.SerialNumber)
		if cert.Id != expected {
			err = fmt.Errorf("expected certificate serial to be %s, was actually %s", expected, cert.Id)
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		// Update the certreq record with the certificate ID
		certreq.Certificate = cert.Id

		// Append the certificate serial to the vasp
		if err = models.AppendCertID(vasp, cert.Id); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		// Put the certificate record
		var putReply *pb.PutReply
		var certData []byte
		if certData, err = proto.Marshal(cert); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}
		if putReply, err = dbClient.Put(ctx, &pb.PutRequest{
			Namespace: wire.NamespaceCerts,
			Key:       []byte(cert.Id),
			Value:     certData,
		}); err != nil {
			err = fmt.Errorf("could not put certificate %s: %s", cert.Id, err)
			results = multierror.Append(results, err)
			continue vaspLoop
		}
		if !putReply.Success {
			err = fmt.Errorf("could not put certificate %s", cert.Id)
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		// Put the vasp record
		var vaspData []byte
		if vaspData, err = proto.Marshal(vasp); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}
		if putReply, err = dbClient.Put(ctx, &pb.PutRequest{
			Namespace: wire.NamespaceVASPs,
			Key:       []byte(vasp.Id),
			Value:     vaspData,
		}); err != nil {
			err = fmt.Errorf("could not put vasp %s: %s", vasp.Id, err)
			results = multierror.Append(results, err)
			continue vaspLoop
		}
		if !putReply.Success {
			err = fmt.Errorf("could not put vasp %s", vasp.Id)
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		// Put the certreq record
		var certreqData []byte
		if certreqData, err = proto.Marshal(certreq); err != nil {
			results = multierror.Append(results, err)
			continue vaspLoop
		}
		if putReply, err = dbClient.Put(ctx, &pb.PutRequest{
			Namespace: wire.NamespaceCertReqs,
			Key:       []byte(certreq.Id),
			Value:     certreqData,
		}); err != nil {
			err = fmt.Errorf("could not put certreq %s: %s", certreq.Id, err)
			results = multierror.Append(results, err)
			continue vaspLoop
		}
		if !putReply.Success {
			err = fmt.Errorf("could not put certreq %s", certreq.Id)
			results = multierror.Append(results, err)
			continue vaspLoop
		}

		migrated++
		fmt.Printf("migrated certreq %s (%s) to certificate %s\n", certreq.Id, certreq.CommonName, cert.Id)
		if report != nil {
			report.WriteString(fmt.Sprintf("migrated certreq %s (%s) to certificate %s\n", certreq.Id, certreq.CommonName, cert.Id))
			if err = writeJSON(cert, report); err != nil {
				results = multierror.Append(results, err)
				break vaspLoop
			}
			report.WriteString("\n\n")
		}
	}

	fmt.Println("Migrated", migrated, "certificates")
	fmt.Println("Skipped", skipped, "vasps")
	if results != nil {
		fmt.Println(results.Error())
	}
	if report != nil {
		fmt.Println("Report written to", c.String("report"))
	}
	return nil
}

func restoreCerts(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	// Open the report file if requested
	var report *os.File
	if c.String("report") != "" {
		if report, err = os.Create(c.String("report")); err != nil {
			return cli.Exit(err, 1)
		}
		defer report.Close()
	}

	// Get all the VASPs in the database, indexed by common name
	vasps := make(map[string]*gds.VASP)
	req := &pb.CursorRequest{
		Namespace: wire.NamespaceVASPs,
	}
	var stream pb.Trtl_CursorClient
	if stream, err = dbClient.Cursor(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}
	for {
		var rep *pb.KVPair
		if rep, err = stream.Recv(); err != nil {
			if err != io.EOF {
				return cli.Exit(err, 1)
			}
			break
		}

		vasp := &gds.VASP{}
		if err = proto.Unmarshal(rep.Value, vasp); err != nil {
			return cli.Exit(err, 1)
		}
		if _, ok := vasps[vasp.CommonName]; ok {
			return cli.Exit(fmt.Errorf("duplicate VASP common name in database: %s", vasp.CommonName), 1)
		}
		vasps[vasp.CommonName] = vasp
	}
	fmt.Printf("found %d VASPs in database\n", len(vasps))
	if report != nil {
		report.WriteString(fmt.Sprintf("found %d VASPs in database\n", len(vasps)))
	}

	// Get all the certificate requests in the database, indexed by batch ID
	certReqs := make(map[int64]*models.CertificateRequest)
	req = &pb.CursorRequest{
		Namespace: wire.NamespaceCertReqs,
	}
	if stream, err = dbClient.Cursor(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}
	for {
		var rep *pb.KVPair
		if rep, err = stream.Recv(); err != nil {
			if err != io.EOF {
				return cli.Exit(err, 1)
			}
			break
		}

		certReq := &models.CertificateRequest{}
		if err = proto.Unmarshal(rep.Value, certReq); err != nil {
			return cli.Exit(err, 1)
		}
		if _, ok := certReqs[certReq.BatchId]; ok {
			return cli.Exit(fmt.Errorf("duplicate certificate request batch ID in database: %d", certReq.BatchId), 1)
		}
		certReqs[certReq.BatchId] = certReq
	}
	fmt.Printf("found %d certificate requests in database\n", len(certReqs))
	if report != nil {
		report.WriteString(fmt.Sprintf("found %d certificate requests in database\n", len(certReqs)))
	}

	// Load the audit log from the CSV file
	var file *os.File
	if file, err = os.Open(c.String("path")); err != nil {
		return cli.Exit(err, 1)
	}
	defer file.Close()

	// Read the CSV file
	var records [][]string
	if records, err = csv.NewReader(file).ReadAll(); err != nil {
		return cli.Exit(err, 1)
	}
	if len(records) == 0 {
		return cli.Exit("no records found in CSV file", 1)
	}

	// Parse the CSV header
	rows := make(map[string]int)
	for i, col := range records[0] {
		rows[col] = i
	}

	var migrated, skipped int
	var results *multierror.Error

	migratedCerts := make(map[string]*models.Certificate)
recordLoop:
	for i, record := range records[1:] {
		cert := &models.Certificate{}

		status := record[rows["Status"]]
		switch status {
		case "ISSUED":
			cert.Status = models.CertificateState_ISSUED
		case "REVOKED":
			cert.Status = models.CertificateState_REVOKED
		default:
			fmt.Printf("skipping record %d with status %s\n", i, status)
			if report != nil {
				report.WriteString(fmt.Sprintf("skipping record %d with status %s\n", i, status))
			}
			skipped++
			continue recordLoop
		}

		id := record[rows["Serial Number"]]
		if id == "" {
			results = multierror.Append(results, fmt.Errorf("record %d has no serial number", i))
			continue recordLoop
		}
		cert.Id = id

		// Make sure there are no duplicates
		if _, ok := migratedCerts[id]; ok {
			fmt.Printf("skipping record %d: duplicate certificate id %s\n", i, id)
			if report != nil {
				report.WriteString(fmt.Sprintf("skipping record %d: duplicate certificate id %s\n", i, id))
			}
			skipped++
			continue recordLoop
		}

		// Only put the certificate record if it doesn't already exist
		if _, err = dbClient.Get(ctx, &pb.GetRequest{
			Namespace: wire.NamespaceCerts,
			Key:       []byte(id),
		}); err == nil {
			fmt.Printf("skipping record %d: certificate already exists\n", i)
			if report != nil {
				report.WriteString(fmt.Sprintf("skipping record %d: certificate already exists\n", i))
			}
			skipped++
			continue recordLoop
		}
		if status, ok := grpc.FromError(err); !ok || status.Code() != codes.NotFound {
			results = multierror.Append(results, err)
			continue recordLoop
		}

		// Parse the issued at timestamp
		var issuedAt time.Time
		const layout = "2006-01-02T15:04:05.000000"
		if issuedAt, err = time.Parse(layout, record[rows["Issued At"]]); err != nil {
			results = multierror.Append(results, fmt.Errorf("record %d has invalid issued at time: %s", i, err))
			continue recordLoop
		}
		expiresAt := issuedAt.AddDate(0, 13, 0)
		if expiresAt.After(time.Now()) {
			cert.Status = models.CertificateState_EXPIRED
		}

		// Parse the subject params into the subject and issuer
		subject := &gds.Name{}
		issuer := &gds.Name{
			CommonName: "CipherTrace Issuing CA",
		}
		var profileId string
		if profileId = record[rows["Subject"]]; profileId == "" {
			results = multierror.Append(results, fmt.Errorf("record %d has no subject params", i))
			continue recordLoop
		}
		params := strings.Split(profileId, ",")
		if len(params) == 1 {
			subject.CommonName = params[0]
		} else {
			// Parse the individual profile params
			for _, param := range params {
				kv := strings.Split(param, "=")
				if len(kv) != 2 {
					results = multierror.Append(results, fmt.Errorf("record %d has invalid profile param: %s", i, param))
					continue recordLoop
				}
				switch kv[0] {
				case "CN":
					subject.CommonName = kv[1]
				case "O":
					subject.Organization = []string{kv[1]}
					issuer.Organization = []string{kv[1]}
				case "L":
					subject.Locality = []string{kv[1]}
					issuer.Locality = []string{kv[1]}
				case "ST":
					subject.Province = []string{kv[1]}
					issuer.Province = []string{kv[1]}
				case "C":
					subject.Country = []string{kv[1]}
					issuer.Country = []string{kv[1]}
				default:
					results = multierror.Append(results, fmt.Errorf("record %d has unknown profile param: %s", i, param))
					continue recordLoop
				}
			}
		}

		if subject.CommonName == "" {
			results = multierror.Append(results, fmt.Errorf("could not parse common name from record %d", i))
			continue recordLoop
		}

		// Get the VASP from the common name
		var vasp *gds.VASP
		var ok bool
		if vasp, ok = vasps[subject.CommonName]; !ok {
			results = multierror.Append(results, fmt.Errorf("could not find VASP with common name %s", subject.CommonName))
			continue recordLoop
		}

		// Get the certificate request from the batch ID
		var batchId int64
		if batchId, err = strconv.ParseInt(record[rows["Batch Id"]], 10, 64); err != nil {
			results = multierror.Append(results, fmt.Errorf("record %d has invalid batch id: %s", i, err))
			continue recordLoop
		}
		var certReq *models.CertificateRequest
		if certReq, ok = certReqs[batchId]; !ok {
			cert.Status = models.CertificateState_REVOKED
			results = multierror.Append(results, fmt.Errorf("WARNING: could not find certificate request with batch id %d", batchId))
		}

		// Amend the records with the reference IDs
		cert.Vasp = vasp.Id
		if certReq != nil {
			cert.Request = certReq.Id
			certReq.Certificate = cert.Id
		}
		if err = models.AppendCertID(vasp, cert.Id); err != nil {
			results = multierror.Append(results, err)
			continue recordLoop
		}

		// Add the certificate details
		cert.Details = &gds.Certificate{
			SerialNumber: []byte(id),
			Subject:      subject,
			Issuer:       issuer,
			NotBefore:    issuedAt.Format(time.RFC3339),
			NotAfter:     expiresAt.Format(time.RFC3339),
			Revoked:      cert.Status == models.CertificateState_REVOKED,
		}

		// Put the certificate record
		var putReply *pb.PutReply
		var certData []byte
		if certData, err = proto.Marshal(cert); err != nil {
			results = multierror.Append(results, err)
			continue recordLoop
		}
		if putReply, err = dbClient.Put(ctx, &pb.PutRequest{
			Namespace: wire.NamespaceCerts,
			Key:       []byte(cert.Id),
			Value:     certData,
		}); err != nil {
			err = fmt.Errorf("could not put certificate %s: %s", cert.Id, err)
			results = multierror.Append(results, err)
			continue recordLoop
		}
		if !putReply.Success {
			err = fmt.Errorf("could not put certificate %s", cert.Id)
			results = multierror.Append(results, err)
			continue recordLoop
		}

		// Put the vasp record
		var vaspData []byte
		if vaspData, err = proto.Marshal(vasp); err != nil {
			results = multierror.Append(results, err)
			continue recordLoop
		}
		if putReply, err = dbClient.Put(ctx, &pb.PutRequest{
			Namespace: wire.NamespaceVASPs,
			Key:       []byte(vasp.Id),
			Value:     vaspData,
		}); err != nil {
			err = fmt.Errorf("could not put vasp %s: %s", vasp.Id, err)
			results = multierror.Append(results, err)
			continue recordLoop
		}
		if !putReply.Success {
			err = fmt.Errorf("could not put vasp %s", vasp.Id)
			results = multierror.Append(results, err)
			continue recordLoop
		}

		if certReq != nil {
			// Put the certreq record
			var certreqData []byte
			if certreqData, err = proto.Marshal(certReq); err != nil {
				results = multierror.Append(results, err)
				continue recordLoop
			}
			if putReply, err = dbClient.Put(ctx, &pb.PutRequest{
				Namespace: wire.NamespaceCertReqs,
				Key:       []byte(certReq.Id),
				Value:     certreqData,
			}); err != nil {
				err = fmt.Errorf("could not put certreq %s: %s", certReq.Id, err)
				results = multierror.Append(results, err)
				continue recordLoop
			}
			if !putReply.Success {
				err = fmt.Errorf("could not put certreq %s", certReq.Id)
				results = multierror.Append(results, err)
				continue recordLoop
			}
		}

		migratedCerts[cert.Id] = cert
		migrated++
		fmt.Printf("migrated certificate %s (%s)\n", cert.Id, cert.Details.Subject.CommonName)
		if report != nil {
			report.WriteString(fmt.Sprintf("migrated certificate %s (%s)\n", cert.Id, cert.Details.Subject.CommonName))
			if err = writeJSON(cert, report); err != nil {
				results = multierror.Append(results, err)
				continue recordLoop
			}
			report.WriteString("\n\n")
		}
	}

	fmt.Println("Migrated", migrated, "records")
	fmt.Println("Skipped", skipped, "records")
	if results != nil {
		fmt.Println(results.Error())
	}
	if report != nil {
		fmt.Println("Report written to", c.String("report"))
	}
	return nil
}

// dbGet prints values from the trtl database given a set of keys.
func dbGet(c *cli.Context) (err error) {
	b64decode := c.Bool("b64decode")
	ctx, cancel := profile.Context()
	defer cancel()

	for _, keys := range c.Args().Slice() {
		var key []byte
		if key, err = wire.DecodeKey(keys, b64decode); err != nil {
			return cli.Exit(fmt.Errorf("could not decode key: %s", err), 1)
		}

		// Execute the Get request
		var resp *pb.GetReply
		req := &pb.GetRequest{
			Key:       key,
			Namespace: c.String("namespace"),
			Options: &pb.Options{
				ReturnMeta: c.Bool("meta"),
			},
		}
		if resp, err = dbClient.Get(ctx, req); err != nil {
			return cli.Exit(err, 1)
		}

		// Print the response using the protojson printer and add a newline after.
		if err = printJSON(resp); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Println("")
	}
	return nil
}

// dbPut puts a value to to a key in the trtl database
func dbPut(c *cli.Context) (err error) {
	// Create the request context
	ctx, cancel := profile.Context()
	defer cancel()

	// Execute the Put request
	var resp *pb.PutReply
	req := &pb.PutRequest{
		Key:       []byte(c.String("key")),
		Value:     []byte(c.String("value")),
		Namespace: c.String("namespace"),
		Options: &pb.Options{
			ReturnMeta: c.Bool("meta"),
		},
	}
	if resp, err = dbClient.Put(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}
	if resp.Success {
		fmt.Printf("successfully put value %s to key %s\n", req.Value, req.Key)
	} else {
		fmt.Printf("could not put value %s to key %s\n", req.Value, req.Key)
	}
	if resp.Meta != nil {
		// Print the response
		if err = printJSON(resp); err != nil {
			return cli.Exit(err, 1)
		}
	}
	return nil
}

// dbDelete deletes a key in the trtl database
func dbDelete(c *cli.Context) (err error) {
	// Create the request context
	ctx, cancel := profile.Context()
	defer cancel()

	// Execute the Put request
	var resp *pb.DeleteReply
	req := &pb.DeleteRequest{
		Key:       []byte(c.String("key")),
		Namespace: c.String("namespace"),
		Options: &pb.Options{
			ReturnMeta: c.Bool("meta"),
		},
	}
	if resp, err = dbClient.Delete(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}
	if resp.Success {
		fmt.Printf("successfully deleted key %s\n", req.Key)
	} else {
		fmt.Printf("could not delete key %s\n", req.Key)
	}
	if resp.Meta != nil {
		// Print the protocol buffer response as JSON
		if err = printJSON(resp); err != nil {
			return cli.Exit(err, 1)
		}
	}

	return nil
}

// dbList lists all the keys in the trtl database
func dbList(c *cli.Context) (err error) {
	// Create the request context
	ctx, cancel := profile.Context()
	defer cancel()

	// Create a cursor request that returns no values, just keys in the specified namespace
	req := &pb.CursorRequest{
		Namespace: c.String("namespace"),
		Options: &pb.Options{
			IterNoValues: true,
		},
	}

	b64decode := c.Bool("b64decode")
	if prefix := c.String("prefix"); prefix != "" {
		if req.Prefix, err = wire.DecodeKey(prefix, b64decode); err != nil {
			return cli.Exit(fmt.Errorf("could not decode prefix: %s", err), 1)
		}
	}

	if seekKey := c.String("seek-key"); seekKey != "" {
		if req.SeekKey, err = wire.DecodeKey(seekKey, b64decode); err != nil {
			return cli.Exit(fmt.Errorf("could not decode seek-key: %s", err), 1)
		}
	}

	var stream pb.Trtl_CursorClient
	if stream, err = dbClient.Cursor(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}

	for {
		var rep *pb.KVPair
		if rep, err = stream.Recv(); err != nil {
			if err != io.EOF {
				return cli.Exit(err, 1)
			}
			break
		}

		fmt.Println(wire.EncodeKey(rep.Key, b64decode))
	}
	return nil
}

// status prints the status of the trtl service.
func status(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	var rep *pb.ServerStatus
	if rep, err = dbClient.Status(ctx, &pb.HealthCheck{}); err != nil {
		return cli.Exit(err, 1)
	}
	return printJSON(rep)
}

//===========================================================================
// Peers (Replica) Client Functions
//===========================================================================

// addPeers creates a Peer and calls the peers management service to add it.
func addPeers(c *cli.Context) (err error) {
	// Create the Peer with the specified PID
	// TODO: how to add the other values for a Peer?
	peer := &peers.Peer{
		Id:     c.Uint64("pid"),
		Addr:   c.String("addr"),
		Name:   c.String("name"),
		Region: c.String("region"),
	}

	if peer.Id == 0 || peer.Addr == "" {
		return cli.Exit("must specify ID and addr to add a peer", 1)
	}

	// create a new context and pass the parent context in
	ctx, cancel := profile.Context()
	defer cancel()

	// call client.AddPeer with the pid
	var out *peers.PeersStatus
	if out, err = peersClient.AddPeers(ctx, peer); err != nil {
		return cli.Exit(err, 1)
	}

	// print the returned result
	printJSON(out)
	return nil
}

// removePeers calls the peers management service to remove a peer.
func delPeers(c *cli.Context) (err error) {
	peer := &peers.Peer{
		Id: c.Uint64("pid"),
	}

	// create a new context and pass the parent context in
	ctx, cancel := profile.Context()
	defer cancel()

	// call client.RmPeer with the pid
	var out *peers.PeersStatus
	if out, err = peersClient.RmPeers(ctx, peer); err != nil {
		return cli.Exit(err, 1)
	}

	// print the returned result
	printJSON(out)
	return nil
}

// listPeers calls the peers management service to list requested peers.
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
	if out, err = peersClient.GetPeers(ctx, filter); err != nil {
		return cli.Exit(err, 1)
	}

	// print the peers
	printJSON(out)
	return nil
}

//===========================================================================
// Anti-Entropy (Replica) Admin Functions
//===========================================================================

func gossip(c *cli.Context) (err error) {
	return errors.New("honu replication required")
}

func gossipMigrate(c *cli.Context) (err error) {
	return errors.New("honu object migration required")
}

//===========================================================================
// Profile Functions
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

	// Handle edit and then exit
	if c.Bool("edit") {
		if err = profiles.EditProfiles(); err != nil {
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
func printJSON(m interface{}) (err error) {
	var data []byte

	switch msg := m.(type) {
	case proto.Message:
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}
		if data, err = opts.Marshal(msg); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		if data, err = json.MarshalIndent(m, "", "  "); err != nil {
			return cli.Exit(err, 1)
		}
	}

	fmt.Println(string(data))
	return nil
}

// helper function to write JSON response to file and exit
func writeJSON(m interface{}, file *os.File) (err error) {
	var data []byte

	switch msg := m.(type) {
	case proto.Message:
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}
		if data, err = opts.Marshal(msg); err != nil {
			return cli.Exit(err, 1)
		}
	default:
		if data, err = json.MarshalIndent(m, "", "  "); err != nil {
			return cli.Exit(err, 1)
		}
	}

	if _, err = file.Write(data); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}
