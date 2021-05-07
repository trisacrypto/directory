package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/trisacrypto/directory/pkg"
	trisads "github.com/trisacrypto/directory/pkg/gds"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	client      api.TRISADirectoryClient
	adminClient admin.DirectoryAdministrationClient
)

func main() {
	app := cli.NewApp()

	app.Name = "trisads"
	app.Version = pkg.Version()
	app.Usage = "a gRPC based directory service for TRISA identity lookups"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "e, endpoint",
			Usage:  "the url to connect the directory service client",
			Value:  "api.vaspdirectory.net:443",
			EnvVar: "TRISA_DIRECTORY_URL",
		},
		cli.BoolFlag{
			Name:  "S, no-secure",
			Usage: "do not connect via TLS (e.g. for development)",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:     "serve",
			Usage:    "run the trisa directory service",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "a, addr",
					Usage:  "the address and port to bind the server on",
					EnvVar: "TRISADS_BIND_ADDR",
				},
			},
		},
		{
			Name:      "load",
			Usage:     "load the directory from a csv file",
			Category:  "server",
			Action:    load,
			ArgsUsage: "csv [csv ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "TRISADS_DATABASE",
				},
			},
		},
		{
			Name:     "review",
			Usage:    "submit a VASP registration review response",
			Category: "admin",
			Action:   review,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i, id",
					Usage: "the ID of the VASP to submit the review for",
				},
				cli.StringFlag{
					Name:  "t, token",
					Usage: "the administrative token sent in the review request email",
				},
				cli.BoolFlag{
					Name:  "R, reject",
					Usage: "reject the registration request",
				},
				cli.BoolFlag{
					Name:  "a, accept",
					Usage: "accept the registration request",
				},
				cli.StringFlag{
					Name:  "m, reason",
					Usage: "provide a reason to reject the request",
				},
			},
		},
		{
			Name:     "register",
			Usage:    "register a VASP using json data",
			Category: "client",
			Action:   register,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "d, data",
					Usage: "the json file containing the VASP data record",
				},
			},
		},
		{
			Name:     "lookup",
			Usage:    "lookup VASPs using name or ID",
			Category: "client",
			Action:   lookup,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i, id",
					Usage: "id of the VASP to lookup",
				},
				cli.StringFlag{
					Name:  "d, directory",
					Usage: "directory that registered the VASP (assumes target directory by default)",
				},
				cli.StringFlag{
					Name:  "n, common-name",
					Usage: "domain name of the VASP to lookup (case-insensitive, exact match)",
				},
			},
		},
		{
			Name:     "search",
			Usage:    "search for VASPs using name or country",
			Category: "client",
			Action:   search,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "n, name",
					Usage: "one or more names of VASPs to search for",
				},
				cli.StringSliceFlag{
					Name:  "w, web",
					Usage: "one or more websites of VASPs to search for",
				},
				cli.StringSliceFlag{
					Name:  "c, country",
					Usage: "one or more countries to filter on",
				},
				cli.StringSliceFlag{
					Name:  "C, category",
					Usage: "one or more categories to filter on",
				},
			},
		},
		{
			Name:     "status",
			Usage:    "check on the verification and service status of a VASP",
			Category: "client",
			Action:   status,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i, id",
					Usage: "id of the VASP to lookup",
				},
				cli.StringFlag{
					Name:  "d, directory",
					Usage: "directory that registered the VASP (assumes target directory by default)",
				},
				cli.StringFlag{
					Name:  "n, common-name",
					Usage: "domain name of the VASP to lookup (case-insensitive, exact match)",
				},
			},
		},
		{
			Name:     "verify",
			Usage:    "verify your email address with the token",
			Category: "client",
			Action:   verifyEmail,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i, id",
					Usage: "id of the VASP your contact information is attached to",
				},
				cli.StringFlag{
					Name:  "t, token",
					Usage: "token that was emailed to you for verification",
				},
			},
		},
	}

	app.Run(os.Args)
}

// Serve the TRISA directory service
func serve(c *cli.Context) (err error) {
	var conf *trisads.Settings
	if conf, err = trisads.Config(); err != nil {
		return cli.NewExitError(err, 1)
	}

	if addr := c.String("addr"); addr != "" {
		conf.BindAddr = addr
	}

	var srv *trisads.Server
	if srv, err = trisads.New(conf); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

// Load the LevelDB database with initial directory info from CSV
// TODO: remove or make more robust
func load(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.NewExitError("specify path to csv data to load", 1)
	}

	var dsn string
	if dsn = c.String("db"); dsn == "" {
		return cli.NewExitError("please specify a dsn to connect to the directory store", 1)
	}

	var db store.Store
	if db, err = store.Open(dsn); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer db.Close()

	for _, path := range c.Args() {
		if err = store.Load(db, path); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

// Submit a review for a registration request
func review(c *cli.Context) (err error) {
	if (!c.Bool("accept") && !c.Bool("reject")) || (c.Bool("accept") && c.Bool("reject")) {
		return cli.NewExitError("specify either accept or reject", 1)
	}

	req := &admin.ReviewRequest{
		Id:                     c.String("id"),
		AdminVerificationToken: c.String("token"),
		Accept:                 c.Bool("accept") && !c.Bool("reject"),
		RejectReason:           c.String("reason"),
	}

	if req.Id == "" || req.AdminVerificationToken == "" {
		return cli.NewExitError("specify both id and token", 1)
	}

	if !req.Accept && req.RejectReason == "" {
		return cli.NewExitError("must specify a reject reason if rejecting", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := adminClient.Review(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// Register an entity using the API from a CLI client
func register(c *cli.Context) (err error) {
	var path string
	if path = c.String("data"); path == "" {
		return cli.NewExitError("specify a json file to load the entity data from", 1)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	var req *api.RegisterRequest
	if err = json.Unmarshal(data, &req); err != nil {
		return cli.NewExitError(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Register(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// Lookup VASPs using the API from a CLI client
func lookup(c *cli.Context) (err error) {
	id := c.String("id")
	directory := c.String("directory")
	commonName := c.String("common-name")

	if commonName == "" && id == "" {
		return cli.NewExitError("specify either name or id for lookup", 1)
	}

	if commonName != "" && id != "" {
		return cli.NewExitError("specify either name or id for lookup, not both", 1)
	}

	req := &api.LookupRequest{
		Id:                  id,
		RegisteredDirectory: directory,
		CommonName:          commonName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Lookup(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// Search for VASPs by name or country using the API from a CLI client
func search(c *cli.Context) (err error) {
	req := &api.SearchRequest{
		Name:             c.StringSlice("name"),
		Website:          c.StringSlice("web"),
		Country:          c.StringSlice("country"),
		BusinessCategory: make([]models.BusinessCategory, 0, len(c.StringSlice("category"))),
		VaspCategory:     make([]models.VASPCategory, 0, len(c.StringSlice("category"))),
	}

	for _, cat := range c.StringSlice("category") {
		if enum, ok := models.BusinessCategory_value[cat]; ok {
			req.BusinessCategory = append(req.BusinessCategory, models.BusinessCategory(enum))
			continue
		}
		if enum, ok := models.VASPCategory_value[cat]; ok {
			req.VaspCategory = append(req.VaspCategory, models.VASPCategory(enum))
			continue
		}
		return cli.NewExitError(fmt.Errorf("unknown category %q", cat), 1)
	}

	if len(req.Name) == 0 && len(req.Website) == 0 {
		return cli.NewExitError("specify search query", 1)
	}

	for _, web := range req.Website {
		if u, err := url.Parse(web); err != nil {
			return cli.NewExitError(fmt.Errorf("%q not a valid URL: %s", web, err), 1)
		} else if u.Hostname() == "" {
			return cli.NewExitError(fmt.Errorf("%q not a valid URL: requires scheme e.g. http://", web), 1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Search(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// Check on verification and service status of a VASP
func status(c *cli.Context) (err error) {
	req := &api.StatusRequest{
		Id:                  c.String("id"),
		RegisteredDirectory: c.String("directory"),
		CommonName:          c.String("common-name"),
	}

	if req.Id == "" && req.CommonName == "" {
		return cli.NewExitError("specify either id or common-name", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Status(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// Send email verification code to the directory serivce
func verifyEmail(c *cli.Context) (err error) {
	req := &api.VerifyEmailRequest{
		Id:    c.String("id"),
		Token: c.String("token"),
	}

	if req.Id == "" || req.Token == "" {
		return cli.NewExitError("specify both id and token", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.VerifyEmail(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// helper function to create the GRPC client with default options
func initClient(c *cli.Context) (err error) {
	var opts []grpc.DialOption
	if c.GlobalBool("no-secure") {
		opts = append(opts, grpc.WithInsecure())
	} else {
		config := &tls.Config{}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.GlobalString("endpoint"), opts...); err != nil {
		return cli.NewExitError(err, 1)
	}
	client = api.NewTRISADirectoryClient(cc)
	adminClient = admin.NewDirectoryAdministrationClient(cc)
	return nil
}

// helper function to print JSON response and exit
func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
