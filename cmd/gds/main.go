package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	profiles "github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/config"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var (
	profile       *profiles.Profile
	client        api.TRISADirectoryClient
	adminClient   admin.DirectoryAdministrationClient
	membersClient members.TRISAMembersClient
)

// Format for YYYY-MM-DD time representation
const weekFormat = "2006-01-02"

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	// Create the CLI application
	app := &cli.App{
		Name:    "gds",
		Version: pkg.Version(),
		Usage:   "the TRISA global directory service",
		Before:  loadProfile,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"P"},
				Usage:   "specify the client profile to use (default: active profile)",
				EnvVars: []string{"TRISA_DIRECTORY_PROFILE", "GDS_PROFILE"},
			},
			&cli.StringFlag{
				Name:    "directory-endpoint",
				Aliases: []string{"e"},
				Usage:   "the url to connect the directory service client",
				EnvVars: []string{"TRISA_DIRECTORY_URL", "GDS_DIRECTORY_URL"},
			},
			&cli.StringFlag{
				Name:    "admin-endpoint",
				Aliases: []string{"a"},
				Usage:   "the url to connect the directory administration client",
				EnvVars: []string{"TRISA_DIRECTORY_ADMIN_URL", "GDS_ADMIN_URL"},
			},
			&cli.StringFlag{
				Name:    "members-endpoint",
				Aliases: []string{"m"},
				Usage:   "the url to connect the trisa members client",
				EnvVars: []string{"TRISA_MEMBERS_URL", "GDS_MEMBERS_URL"},
			},
			&cli.BoolFlag{
				Name:    "no-secure",
				Aliases: []string{"S"},
				Usage:   "do not connect via TLS (e.g. for development)",
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "serve",
				Usage:    "run the trisa directory service",
				Category: "server",
				Action:   serve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "addr",
						Aliases: []string{"a"},
						Usage:   "the address and port to bind the server on",
						EnvVars: []string{"GDS_BIND_ADDR"},
					},
				},
			},
			{
				// TODO: move this to gdsutil as it is deprecated
				Name:      "load",
				Usage:     "load the directory from a csv file",
				Category:  "server",
				Action:    load,
				ArgsUsage: "csv [csv ...]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "db",
						Aliases: []string{"d"},
						Usage:   "dsn to connect to gds directory storage",
						EnvVars: []string{"GDS_DATABASE_URL"},
					},
				},
			},
			{
				Name:     "gds:register",
				Usage:    "register a VASP using json data",
				Category: "gds",
				Action:   register,
				Before:   initClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "data",
						Aliases: []string{"d"},
						Usage:   "the JSON file containing the VASP data record",
					},
				},
			},
			{
				Name:     "gds:lookup",
				Usage:    "lookup VASPs using name or ID",
				Category: "gds",
				Action:   lookup,
				Before:   initClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "id of the VASP to lookup",
					},
					&cli.StringFlag{
						Name:    "directory",
						Aliases: []string{"d"},
						Usage:   "directory that registered the VASP (assumes target directory by default)",
					},
					&cli.StringFlag{
						Name:    "common-name",
						Aliases: []string{"n"},
						Usage:   "domain name of the VASP to lookup (case-insensitive, exact match)",
					},
				},
			},
			{
				Name:     "gds:search",
				Usage:    "search for VASPs using name or country",
				Category: "gds",
				Action:   search,
				Before:   initClient,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "one or more names of VASPs to search for",
					},
					&cli.StringSliceFlag{
						Name:    "web",
						Aliases: []string{"w"},
						Usage:   "one or more websites of VASPs to search for",
					},
					&cli.StringSliceFlag{
						Name:    "country",
						Aliases: []string{"c"},
						Usage:   "one or more countries to filter on",
					},
					&cli.StringSliceFlag{
						Name:    "category",
						Aliases: []string{"C", "categories", "cats"},
						Usage:   "one or more categories to filter on",
					},
				},
			},
			{
				Name:     "gds:verification",
				Usage:    "check on the verification and service status of a VASP",
				Category: "gds",
				Action:   verification,
				Before:   initClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "id of the VASP to lookup",
					},
					&cli.StringFlag{
						Name:    "directory",
						Aliases: []string{"d"},
						Usage:   "directory that registered the VASP (assumes target directory by default)",
					},
					&cli.StringFlag{
						Name:    "common-name",
						Aliases: []string{"n"},
						Usage:   "domain name of the VASP to lookup (case-insensitive, exact match)",
					},
				},
			},
			{
				Name:     "gds:verify-contact",
				Usage:    "verify your email address with the token",
				Category: "gds",
				Action:   verifyContact,
				Before:   initClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "id of the VASP your contact information is attached to",
					},
					&cli.StringFlag{
						Name:    "token",
						Aliases: []string{"t"},
						Usage:   "token that was emailed to you for verification",
					},
				},
			},
			{
				Name:     "gds:status",
				Usage:    "send a health check request to the directory service",
				Category: "gds",
				Action:   status,
				Before:   initClient,
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "attempts",
						Aliases: []string{"a"},
						Usage:   "set the number of previous attempts",
					},
					&cli.DurationFlag{
						Name:    "last-checked",
						Aliases: []string{"l"},
						Usage:   "set the last checked field as this long ago",
					},
				},
			},
			{
				Name:     "admin:review",
				Usage:    "submit a VASP registration review response",
				Category: "admin",
				Action:   review,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the ID of the VASP to submit the review for",
					},
					&cli.StringFlag{
						Name:    "token",
						Aliases: []string{"t"},
						Usage:   "the administrative token sent in the review request email",
					},
					&cli.BoolFlag{
						Name:    "fetch-token",
						Aliases: []string{"T"},
						Usage:   "attempt to fetch the admin verification token before review",
					},
					&cli.BoolFlag{
						Name:    "reject",
						Aliases: []string{"R"},
						Usage:   "reject the registration request",
					},
					&cli.BoolFlag{
						Name:    "accept",
						Aliases: []string{"a"},
						Usage:   "accept the registration request",
					},
					&cli.StringFlag{
						Name:    "reason",
						Aliases: []string{"m"},
						Usage:   "provide a reason to reject the request",
					},
				},
			},
			{
				Name:     "admin:resend",
				Usage:    "request emails be resent in case of delivery errors",
				Category: "admin",
				Action:   resend,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the ID of the VASP to submit the review for",
					},
					&cli.BoolFlag{
						Name:    "verify-contact",
						Aliases: []string{"v"},
						Usage:   "resend verify contact emails",
					},
					&cli.BoolFlag{
						Name:    "review",
						Aliases: []string{"r"},
						Usage:   "resend review request emails",
					},
					&cli.BoolFlag{
						Name:    "deliver-certs",
						Aliases: []string{"d"},
						Usage:   "resend certificate delivery email",
					},
					&cli.BoolFlag{
						Name:    "reject",
						Aliases: []string{"R"},
						Usage:   "resend rejection email",
					},
					&cli.StringFlag{
						Name:    "reason",
						Aliases: []string{"m"},
						Usage:   "provide a reason to reject the request",
					},
				},
			},
			{
				Name:     "admin:reviews",
				Usage:    "request a timeline of VASP state changes",
				Category: "admin",
				Action:   adminReviewTimeline,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "start",
						Aliases: []string{"s"},
						Usage:   "start date (YYYY-MM-DD) for the review timeline",
						Value:   time.Now().AddDate(-1, 0, 0).Format(weekFormat),
					},
					&cli.StringFlag{
						Name:    "end",
						Aliases: []string{"e"},
						Usage:   "end date (YYYY-MM-DD) for the review timeline",
						Value:   time.Now().Format(weekFormat),
					},
				},
			},
			{
				Name:     "admin:status",
				Usage:    "perform a health check against the admin API",
				Category: "admin",
				Action:   adminStatus,
				Before:   initAdminClient,
				Flags:    []cli.Flag{},
			},
			{
				Name:     "admin:summary",
				Usage:    "collect aggregate information about current GDS status",
				Category: "admin",
				Action:   adminSummary,
				Before:   initAdminClient,
				Flags:    []cli.Flag{},
			},
			{
				Name:     "admin:autocomplete",
				Usage:    "get autocomplete names for the admin searchbar",
				Category: "admin",
				Action:   adminAutocomplete,
				Before:   initAdminClient,
				Flags:    []cli.Flag{},
			},
			{
				Name:     "admin:list",
				Usage:    "list all VASPs summary detail",
				Category: "admin",
				Action:   adminListVASPs,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "page",
						Aliases: []string{"p"},
						Usage:   "query for the specific page of vasps",
					},
					&cli.IntFlag{
						Name:    "page-size",
						Aliases: []string{"s"},
						Usage:   "specify the number of items per page",
					},
					&cli.StringSliceFlag{
						Name:    "status-filters",
						Aliases: []string{"S"},
						Usage:   "filter by verification status",
					},
				},
			},
			{
				Name:     "admin:detail",
				Usage:    "retrieve a VASP detail record by id",
				Category: "admin",
				Action:   adminRetrieveVASP,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the uuid of the VASP to retrieve",
					},
				},
			},
			{
				Name:     "admin:update",
				Usage:    "update a VASP detail record by id",
				Category: "admin",
				Action:   adminUpdateVASP,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the uuid of the VASP to retrieve",
					},
					&cli.StringFlag{
						Name:    "data",
						Aliases: []string{"d"},
						Usage:   "path to JSON data to PATCH VASP detail",
					},
				},
			},
			{
				Name:     "admin:notes",
				Usage:    "list notes associated with a VASP",
				Category: "admin",
				Action:   adminListNotes,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the uuid of the VASP to add a note for",
					},
				},
			},
			{
				Name:     "admin:notes-create",
				Usage:    "create a new note associated with a VASP",
				Category: "admin",
				Action:   adminCreateNote,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the uuid of the VASP to associate the note with",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "the name of the new note or existing note",
					},
					&cli.StringFlag{
						Name:    "text",
						Aliases: []string{"t"},
						Usage:   "the text to include in the note",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "read note text from file",
					},
				},
			},
			{
				Name:     "admin:notes-update",
				Usage:    "update an existing VASP note",
				Category: "admin",
				Action:   adminUpdateNote,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the uuid of the VASP the note is associated with",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "the name of the note to update",
					},
					&cli.StringFlag{
						Name:    "text",
						Aliases: []string{"t"},
						Usage:   "the text to include in the note",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "read note text from file",
					},
				},
			},
			{
				Name:     "admin:notes-delete",
				Usage:    "delete an existing VASP note",
				Category: "admin",
				Action:   adminDeleteNote,
				Before:   initAdminClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "id",
						Aliases: []string{"i"},
						Usage:   "the uuid of the VASP the note is associated with",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "the name of the note to delete",
					},
				},
			},
			{
				Name:     "members:list",
				Usage:    "list all currently verified VASPs in the directory",
				Category: "members",
				Action:   membersList,
				Before:   initMembersClient,
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "page-size",
						Aliases: []string{"s", "size"},
						Usage:   "the number of results per page",
					},
					&cli.StringFlag{
						Name:    "page-token",
						Aliases: []string{"t", "token"},
						Usage:   "next page token for follow-on requests",
					},
					&cli.BoolFlag{
						Name:    "fetc-all",
						Aliases: []string{"a", "all"},
						Usage:   "keep fetching results as long as a next page token is returned",
					},
				},
			},
			{
				Name:      "profile",
				Aliases:   []string{"config"},
				Usage:     "view and manage profiles to configure client with",
				UsageText: "gds profile [name]\n   gds profile --activate [name]\n   gds profile --list\n   gds profile --path\n   gds profile --install",
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
					&cli.StringFlag{
						Name:    "activate",
						Aliases: []string{"a"},
						Usage:   "activate the profile with the specified name",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// CLI Actions
//===========================================================================

// Serve the TRISA directory service
func serve(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if addr := c.String("addr"); addr != "" {
		conf.GDS.BindAddr = addr
	}

	var srv *gds.Service
	if srv, err = gds.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

// Load the LevelDB database with initial directory info from CSV
// TODO: remove or make more robust
func load(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.Exit("specify path to csv data to load", 1)
	}

	var dsn string
	if dsn = c.String("db"); dsn == "" {
		return cli.Exit("please specify a dsn to connect to the directory store", 1)
	}

	var db store.Store
	if db, err = store.Open(config.DatabaseConfig{URL: dsn}); err != nil {
		return cli.Exit(err, 1)
	}
	defer db.Close()

	for _, path := range c.Args().Slice() {
		if err = store.Load(db, path); err != nil {
			return cli.Exit(err, 1)
		}
	}

	return nil
}

// Submit a review for a registration request
func review(c *cli.Context) (err error) {
	if (!c.Bool("accept") && !c.Bool("reject")) || (c.Bool("accept") && c.Bool("reject")) {
		return cli.Exit("specify either accept or reject", 1)
	}

	req := &admin.ReviewRequest{
		ID:                     c.String("id"),
		AdminVerificationToken: c.String("token"),
		Accept:                 c.Bool("accept") && !c.Bool("reject"),
		RejectReason:           c.String("reason"),
	}

	if req.ID == "" {
		return cli.Exit("must specify the id of the VASP", 1)
	}

	if !req.Accept && req.RejectReason == "" {
		return cli.Exit("must specify a reject reason if rejecting", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if c.Bool("fetch-token") {
		rep, err := adminClient.ReviewToken(ctx, req.ID)
		if err != nil {
			return cli.Exit(err, 1)
		}
		req.AdminVerificationToken = rep.AdminVerificationToken
		fmt.Printf("admin verification token fetch: %q\n", rep.AdminVerificationToken)
	}

	if req.AdminVerificationToken == "" {
		return cli.Exit("must specify fetch-token or the admin verification token", 1)
	}

	rep, err := adminClient.Review(ctx, req)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

// Register an entity using the API from a CLI client
func register(c *cli.Context) (err error) {
	var path string
	if path = c.String("data"); path == "" {
		return cli.Exit("specify a json file to load the entity data from", 1)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cli.Exit(err, 1)
	}

	// Check if this is a form downloaded from the UI
	tmp := make(map[string]interface{})
	if err = json.Unmarshal(data, &tmp); err == nil {
		if form, ok := tmp["registrationForm"]; ok {
			if data, err = json.Marshal(form); err != nil {
				return cli.Exit(fmt.Errorf("could not extract registration form: %s", err), 1)
			}
		}
	}

	var req *api.RegisterRequest
	if err = json.Unmarshal(data, &req); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Register(ctx, req)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

// Lookup VASPs using the API from a CLI client
func lookup(c *cli.Context) (err error) {
	id := c.String("id")
	directory := c.String("directory")
	commonName := c.String("common-name")

	if commonName == "" && id == "" {
		return cli.Exit("specify either name or id for lookup", 1)
	}

	if commonName != "" && id != "" {
		return cli.Exit("specify either name or id for lookup, not both", 1)
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
		return cli.Exit(err, 1)
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
		return cli.Exit("specify search query", 1)
	}

	for _, web := range req.Website {
		if u, err := url.Parse(web); err != nil {
			return cli.Exit(fmt.Errorf("%q not a valid URL: %s", web, err), 1)
		} else if u.Hostname() == "" {
			return cli.Exit(fmt.Errorf("%q not a valid URL: requires scheme e.g. http://", web), 1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Search(ctx, req)
	if err != nil {
		return cli.Exit(err, 1)
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
		return cli.Exit("specify either id or common-name", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.Verification(ctx, req)
	if err != nil {
		return cli.Exit(err, 1)
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
		return cli.Exit("specify both id and token", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rep, err := client.VerifyContact(ctx, req)
	if err != nil {
		return cli.Exit(err, 1)
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

	ctx, cancel := profile.Context()
	defer cancel()

	var rep *api.ServiceState
	if rep, err = client.Status(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func resend(c *cli.Context) (err error) {
	req := &admin.ResendRequest{
		ID:     c.String("id"),
		Reason: c.String("reason"),
	}

	if req.ID == "" {
		return cli.Exit("missing VASP record ID, specify with --id", 1)
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
		return cli.Exit("must specify request type (--verify-contact, --review, --deliver-certs, --reject)", 1)
	}

	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.ResendReply
	if rep, err = adminClient.Resend(ctx, req); err != nil {
		return cli.Exit(err, 1)
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
		return cli.Exit(err, 1)
	}
	if _, err = time.Parse(weekFormat, params.End); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.ReviewTimelineReply
	if rep, err = adminClient.ReviewTimeline(ctx, params); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminStatus(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.StatusReply
	if rep, err = adminClient.Status(ctx); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminSummary(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.SummaryReply
	if rep, err = adminClient.Summary(ctx); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminAutocomplete(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.AutocompleteReply
	if rep, err = adminClient.Autocomplete(ctx); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminListVASPs(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	params := &admin.ListVASPsParams{
		Page:          c.Int("page"),
		PageSize:      c.Int("page-size"),
		StatusFilters: c.StringSlice("status-filters"),
	}

	var rep *admin.ListVASPsReply
	if rep, err = adminClient.ListVASPs(ctx, params); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminRetrieveVASP(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.RetrieveVASPReply
	if rep, err = adminClient.RetrieveVASP(ctx, c.String("id")); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminUpdateVASP(c *cli.Context) (err error) {
	req := &admin.UpdateVASPRequest{}
	if path := c.String("data"); path != "" {
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			return cli.Exit(err, 1)
		}

		if err = json.Unmarshal(data, req); err != nil {
			return cli.Exit(fmt.Errorf("could not unmarshal UpdateVASPRequest: %s", err), 1)
		}
	} else {
		return cli.Exit("specify path to JSON data with update vasp request", 1)
	}

	if vaspID := c.String("id"); vaspID != "" {
		req.VASP = vaspID
	}

	ctx, cancel := profile.Context()
	defer cancel()

	var rep *admin.UpdateVASPReply
	if rep, err = adminClient.UpdateVASP(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminListNotes(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	vaspID := c.String("id")
	if vaspID == "" {
		cli.Exit("must specify VASP ID (--id)", 1)
	}

	var rep *admin.ListReviewNotesReply
	if rep, err = adminClient.ListReviewNotes(ctx, vaspID); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminCreateNote(c *cli.Context) (err error) {
	var (
		params *admin.ModifyReviewNoteRequest
		text   string
		file   string
	)

	ctx, cancel := profile.Context()
	defer cancel()

	params = &admin.ModifyReviewNoteRequest{
		VASP:   c.String("id"),
		NoteID: c.String("name"),
	}

	if params.VASP == "" {
		cli.Exit("must specify VASP ID (--id)", 1)
	}

	// Get the note text
	text = c.String("text")
	file = c.String("file")
	if text == "" && file == "" {
		return cli.Exit("must specify either --text or --file", 1)
	} else if text != "" && file != "" {
		return cli.Exit("cannot specify both --text and --file", 1)
	} else if text != "" {
		params.Text = text
	} else {
		var data []byte
		if data, err = os.ReadFile(file); err != nil {
			return cli.Exit(err, 1)
		}
		params.Text = string(data)
	}

	var rep *admin.ReviewNote
	if rep, err = adminClient.CreateReviewNote(ctx, params); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminUpdateNote(c *cli.Context) (err error) {
	var (
		params *admin.ModifyReviewNoteRequest
		text   string
		file   string
	)

	ctx, cancel := profile.Context()
	defer cancel()

	params = &admin.ModifyReviewNoteRequest{
		VASP:   c.String("id"),
		NoteID: c.String("name"),
	}

	if params.VASP == "" {
		cli.Exit("must specify VASP ID (--id)", 1)
	}

	if params.NoteID == "" {
		return cli.Exit("must specify note name (--name)", 1)
	}

	// Get the note text
	text = c.String("text")
	file = c.String("file")
	if text == "" && file == "" {
		return cli.Exit("must specify either --text or --file", 1)
	} else if text != "" && file != "" {
		return cli.Exit("cannot specify both --text and --file", 1)
	} else if text != "" {
		params.Text = text
	} else {
		var data []byte
		if data, err = os.ReadFile(file); err != nil {
			return cli.Exit(err, 1)
		}
		params.Text = string(data)
	}

	var rep *admin.ReviewNote
	if rep, err = adminClient.UpdateReviewNote(ctx, params); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func adminDeleteNote(c *cli.Context) (err error) {
	ctx, cancel := profile.Context()
	defer cancel()

	vaspID := c.String("id")
	if vaspID == "" {
		return cli.Exit("must specify VASP ID (--id)", 1)
	}

	noteID := c.String("name")
	if noteID == "" {
		return cli.Exit("must specify note name (--name)", 1)
	}

	var rep *admin.Reply
	if rep, err = adminClient.DeleteReviewNote(ctx, vaspID, noteID); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

func membersList(c *cli.Context) (err error) {
	// Only fetch a single request if not fetching all
	if !c.Bool("fetch-all") {
		ctx, cancel := profile.Context()
		defer cancel()

		req := &members.ListRequest{
			PageSize:  int32(c.Int64("page-size")),
			PageToken: c.String("page-token"),
		}

		var rep *members.ListReply
		if rep, err = membersClient.List(ctx, req); err != nil {
			return cli.Exit(err, 1)
		}
		return printJSON(rep)
	}

	// Otherwise, keep fetching results until the server has no more
	req := &members.ListRequest{
		PageSize: int32(c.Int64("page-size")),
	}

	for {
		var rep *members.ListReply
		ctx, cancel := profile.Context()
		rep, err = membersClient.List(ctx, req)
		cancel()

		if err != nil {
			return cli.Exit(err, 1)
		}

		for _, member := range rep.Vasps {
			if err = printJSON(member); err != nil {
				return cli.Exit(err, 1)
			}
		}

		if rep.NextPageToken == "" {
			return nil
		}
		req.PageToken = rep.NextPageToken
	}
}

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
// Helper Methods
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

// helper function to create the GRPC client with default options
func initClient(c *cli.Context) (err error) {
	if client, err = profile.Directory.Connect(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func initAdminClient(c *cli.Context) (err error) {
	if adminClient, err = profile.Admin.Connect(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func initMembersClient(c *cli.Context) (err error) {
	if membersClient, err = profile.Members.Connect(); err != nil {
		return cli.Exit(err, 1)
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
			return cli.Exit(err, 1)
		}
	default:
		if data, err = json.MarshalIndent(msg, "", "  "); err != nil {
			return cli.Exit(err, 1)
		}
	}

	fmt.Println(string(data))
	return nil
}
