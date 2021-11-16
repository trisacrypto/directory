package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	profiles "github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/store/wire"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
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
	app.Name = "trtl"
	app.Version = pkg.Version()
	app.Usage = "a command line tool for interacting with the trtl data replication service"
	app.Before = loadProfile
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
					EnvVars: []string{"GDS_DATABASE_URL"},
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
			Name:     "status",
			Usage:    "check the status of the trtl database and replication service",
			Category: "client",
			Before:   initDBClient,
			Action:   status,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect to the trtl replication service",
					EnvVars: []string{"TRISA_DIRECTORY_REPLICA_URL"},
				},
				&cli.BoolFlag{
					Name:    "S",
					Aliases: []string{"no-secure"},
					Usage:   "do not connect via TLS (e.g. for development)",
				},
			},
		},
		{
			Name:      "db:get",
			Usage:     "get a value from the trtl database",
			ArgsUsage: "key [key ...]",
			Category:  "client",
			Before:    initDBClient,
			Action:    dbGet,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect to the trtl replication service",
					EnvVars: []string{"TRISA_DIRECTORY_REPLICA_URL"},
				},
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
				&cli.BoolFlag{
					Name:    "S",
					Aliases: []string{"no-secure"},
					Usage:   "do not connect via TLS (e.g. for development)",
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
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect to the trtl replication service",
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
			Category: "client",
			Before:   initPeersClient,
			Action:   delPeers,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect the trtl replication service",
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
			Category: "client",
			Before:   initPeersClient,
			Action:   listPeers,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "replica-endpoint",
					Aliases: []string{"u"},
					Usage:   "the url to connect the trtl replication service",
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
			Aliases:   []string{"config"},
			Usage:     "view and manage profiles to configure trtl with",
			UsageText: "trtl profile [name]\n   trtl profile --activate [name]\n   trtl profile --list\n   trtl profile --path\n   trtl profile --install",
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
	// TODO: load and validate the trtl configuration
	fmt.Println("trtl config is valid; ready to serve")
	return nil
}

//===========================================================================
// Client Functions
//===========================================================================

var dbClient pb.TrtlClient
var peersClient peers.PeerManagementClient

// initDBClient starts a trtl client with a connection to a trtl database.
func initDBClient(c *cli.Context) (err error) {
	if profile.Replica == nil {
		return cli.Exit("replica not configured", 1)
	}

	if dbClient, err = profile.Replica.ConnectDB(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

// initPeersClient starts a trtl client with a connection to a trtl database.
func initPeersClient(c *cli.Context) (err error) {
	if profile.Replica == nil {
		return cli.Exit("replica not configured", 1)
	}

	if peersClient, err = profile.Replica.ConnectPeers(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

// status prints the status of the trtl service.
func status(c *cli.Context) (err error) {
	// TODO: call the trtl status RPC once implemented
	fmt.Println("trtl status: available")
	return nil
}

// dbGet prints values from the trtl database given a set of keys.
func dbGet(c *cli.Context) (err error) {
	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args().Slice() {
		var key []byte
		if key, err = wire.DecodeKey(keys, b64decode); err != nil {
			return cli.Exit(fmt.Errorf("could not decode key: %s", err), 1)
		}

		var resp *pb.GetReply
		req := &pb.GetRequest{
			Key: key,
			Options: &pb.Options{
				ReturnMeta: c.Bool("meta"),
			},
		}
		if resp, err = dbClient.Get(context.TODO(), req); err != nil {
			return cli.Exit(err, 1)
		}

		jsonpb := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}
		var outdata []byte
		if outdata, err = jsonpb.Marshal(resp); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Println(string(outdata) + "\n")
	}
	return nil
}

// addPeers creates a Peer and calls the peers management service to add it.
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
