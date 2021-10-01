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

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/store"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	client      api.TRISADirectoryClient
	adminClient admin.DirectoryAdministrationClient
)

const weekFormat = "YYYY-MM-DD"

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "gds"
	app.Version = pkg.Version()
	app.Usage = "the global directory service for TRISA"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "e, endpoint",
			Usage:  "the url to connect the directory service client",
			Value:  "api.vaspdirectory.net:443",
			EnvVar: "TRISA_DIRECTORY_URL",
		},
		cli.StringFlag{
			Name:   "a, admin-endpoint",
			Usage:  "the url to connect the directory administration client",
			Value:  "https://api.admin.vaspdirectory.net",
			EnvVar: "TRISA_DIRECTORY_ADMIN_URL",
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
					EnvVar: "GDS_BIND_ADDR",
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
					EnvVar: "GDS_DATABASE_URL",
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
			Name:     "verification",
			Usage:    "check on the verification and service status of a VASP",
			Category: "client",
			Action:   verification,
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
			Name:     "verify-contact",
			Usage:    "verify your email address with the token",
			Category: "client",
			Action:   verifyContact,
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
		{
			Name:     "status",
			Usage:    "send a health check request to the directory service",
			Category: "client",
			Action:   status,
			Before:   initClient,
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:  "a, attempts",
					Usage: "set the number of previous attempts",
				},
				&cli.DurationFlag{
					Name:  "l, last-checked",
					Usage: "set the last checked field as this long ago",
				},
			},
		},
		{
			Name:     "resend",
			Usage:    "request emails be resent in case of delivery errors",
			Category: "admin",
			Action:   resend,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i, id",
					Usage: "the ID of the VASP to submit the review for",
				},
				cli.BoolFlag{
					Name:  "v, verify-contact",
					Usage: "resend verify contact emails",
				},
				cli.BoolFlag{
					Name:  "r, review",
					Usage: "resend review request emails",
				},
				cli.BoolFlag{
					Name:  "d, deliver-certs",
					Usage: "resend certificate delivery email",
				},
				cli.BoolFlag{
					Name:  "R, reject",
					Usage: "resend rejection email",
				},
				cli.StringFlag{
					Name:  "m, reason",
					Usage: "provide a reason to reject the request",
				},
			},
		},
		{
			Name:     "admin:reviews",
			Usage:    "request a timeline of VASP state changes",
			Category: "admin",
			Action:   adminReviewTimeline,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "s, start",
					Usage: "start date (YYYY-MM-DD) for the review timeline",
					Value: time.Now().Add(time.Duration(-52) * time.Hour * 24 * 7).Format(weekFormat),
				},
				cli.StringFlag{
					Name:  "e, end",
					Usage: "end date (YYYY-MM-DD) for the review timeline",
					Value: time.Now().Format(weekFormat),
				},
			},
		},
		{
			Name:     "admin:status",
			Usage:    "perform a health check against the admin API",
			Category: "admin",
			Action:   adminStatus,
			Before:   initClient,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "admin:summary",
			Usage:    "collect aggregate information about current GDS status",
			Category: "admin",
			Action:   adminSummary,
			Before:   initClient,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "admin:autocomplete",
			Usage:    "get autocomplete names for the admin searchbar",
			Category: "admin",
			Action:   adminAutocomplete,
			Before:   initClient,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "admin:list",
			Usage:    "list all VASPs summary detail",
			Category: "admin",
			Action:   adminListVASPs,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "p, page",
					Usage: "query for the specific page of vasps",
				},
				cli.IntFlag{
					Name:  "s, page-size",
					Usage: "specify the number of items per page",
				},
				cli.StringFlag{
					Name:  "S, status",
					Usage: "filter by verification status",
				},
			},
		},
		{
			Name:     "admin:detail",
			Usage:    "retrieve a VASP detail recrod by id",
			Category: "admin",
			Action:   adminRetrieveVASPs,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "i, id",
					Usage: "the uuid of the VASP to retrieve",
				},
			},
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// CLI Actions
//===========================================================================

// Serve the TRISA directory service
func serve(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.NewExitError(err, 1)
	}

	if addr := c.String("addr"); addr != "" {
		conf.GDS.BindAddr = addr
	}

	var srv *gds.Service
	if srv, err = gds.New(conf); err != nil {
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
	if db, err = store.Open(config.DatabaseConfig{URL: dsn}); err != nil {
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
		ID:                     c.String("id"),
		AdminVerificationToken: c.String("token"),
		Accept:                 c.Bool("accept") && !c.Bool("reject"),
		RejectReason:           c.String("reason"),
	}

	if req.ID == "" || req.AdminVerificationToken == "" {
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

	// Check if this is a form downloaded from the UI
	tmp := make(map[string]interface{})
	if err = json.Unmarshal(data, &tmp); err == nil {
		if form, ok := tmp["registrationForm"]; ok {
			if data, err = json.Marshal(form); err != nil {
				return cli.NewExitError(fmt.Errorf("could not extract registration form: %s", err), 1)
			}
		}
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
		VaspCategory:     make([]string, 0, len(c.StringSlice("category"))),
	}

	for _, cat := range c.StringSlice("category") {
		if enum, ok := models.BusinessCategory_value[cat]; ok {
			req.BusinessCategory = append(req.BusinessCategory, models.BusinessCategory(enum))
		} else {
			req.VaspCategory = append(req.VaspCategory, cat)
		}
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
func verification(c *cli.Context) (err error) {
	req := &api.VerificationRequest{
		Id:                  c.String("id"),
		RegisteredDirectory: c.String("directory"),
		CommonName:          c.String("common-name"),
	}

	if req.Id == "" && req.CommonName == "" {
		return cli.NewExitError("specify either id or common-name", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Verification(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// Send email verification code to the directory serivce
func verifyContact(c *cli.Context) (err error) {
	req := &api.VerifyContactRequest{
		Id:    c.String("id"),
		Token: c.String("token"),
	}

	if req.Id == "" || req.Token == "" {
		return cli.NewExitError("specify both id and token", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.VerifyContact(ctx, req)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func status(c *cli.Context) (err error) {
	req := &api.HealthCheck{
		Attempts: uint32(c.Uint("attempts")),
	}

	if lastCheckedAgo := c.Duration("last-checked"); lastCheckedAgo != 0 {
		req.LastCheckedAt = time.Now().Add(-1 * lastCheckedAgo).Format(time.RFC3339)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *api.ServiceState
	if rep, err = client.Status(ctx, req); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func resend(c *cli.Context) (err error) {
	req := &admin.ResendRequest{
		ID:     c.String("id"),
		Reason: c.String("reason"),
	}

	if req.ID == "" {
		return cli.NewExitError("missing VASP record ID, specify with --id", 1)
	}

	// NOTE: if multiple type flags are specified, only one will be used
	switch {
	case c.Bool("verify-contact"):
		req.Action = admin.ResendVerifyContact
	case c.Bool("review"):
		req.Action = admin.ResendReview
	case c.Bool("deliver-certs"):
		req.Action = admin.ResendDeliverCerts
	case c.Bool("reject"):
		req.Action = admin.ResendRejection
	default:
		return cli.NewExitError("must specify request type (--verify-contact, --review, --deliver-certs, --reject)", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *admin.ResendReply
	if rep, err = adminClient.Resend(ctx, req); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func adminReviewTimeline(c *cli.Context) (err error) {
	params := &admin.ReviewTimelineParams{
		Start: c.String("start"),
		End:   c.String("end"),
	}

	// Validate start and end dates
	if _, err = time.Parse(weekFormat, params.Start); err != nil {
		return cli.NewExitError(err, 1)
	}
	if _, err = time.Parse(weekFormat, params.End); err != nil {
		return cli.NewExitError(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *admin.ReviewTimelineReply
	if rep, err = adminClient.ReviewTimeline(ctx, params); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func adminStatus(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *admin.StatusReply
	if rep, err = adminClient.Status(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func adminSummary(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *admin.SummaryReply
	if rep, err = adminClient.Summary(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func adminAutocomplete(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *admin.AutocompleteReply
	if rep, err = adminClient.Autocomplete(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func adminListVASPs(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &admin.ListVASPsParams{
		Page:     c.Int("page"),
		PageSize: c.Int("page-size"),
		Status:   c.String("status"),
	}

	var rep *admin.ListVASPsReply
	if rep, err = adminClient.ListVASPs(ctx, params); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func adminRetrieveVASPs(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *admin.RetrieveVASPReply
	if rep, err = adminClient.RetrieveVASP(ctx, c.String("id")); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

//===========================================================================
// Helper Methods
//===========================================================================

// helper function to create the GRPC client with default options
func initClient(c *cli.Context) (err error) {
	var opts []grpc.DialOption
	if c.GlobalBool("no-secure") {
		opts = append(opts, grpc.WithInsecure())
	} else {
		config := &tls.Config{}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}

	// Connect the directory client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.GlobalString("endpoint"), opts...); err != nil {
		return cli.NewExitError(err, 1)
	}
	client = api.NewTRISADirectoryClient(cc)

	// Connect the admin client
	if adminClient, err = admin.New(c.GlobalString("admin-endpoint")); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

// helper function to print JSON response and exit
func printJSON(msg interface{}) (err error) {

	var data []byte
	switch m := msg.(type) {
	case proto.Message:
		opts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			AllowPartial:    true,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}

		if data, err = opts.Marshal(m); err != nil {
			return cli.NewExitError(err, 1)
		}
	default:
		if data, err = json.MarshalIndent(msg, "", "  "); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	fmt.Println(string(data))
	return nil
}
