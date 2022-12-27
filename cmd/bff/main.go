package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/clive"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	storeconfig "github.com/trisacrypto/directory/pkg/store/config"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"

	"github.com/urfave/cli/v2"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	// Create the CLI application
	app := &cli.App{
		Name:    "gds-bff",
		Version: pkg.Version(),
		Usage:   "a backend for front-end for the GDS service",
		Flags:   []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:     "serve",
				Usage:    "run the gds bff server",
				Category: "server",
				Action:   serve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "addr",
						Aliases: []string{"a"},
						Usage:   "the address and port to bind the server on",
						EnvVars: []string{"GDS_BFF_BIND_ADDR"},
					},
				},
			},
			{
				Name:     "validate",
				Usage:    "validate the current bff configuration",
				Category: "server",
				Action:   validate,
			},
			{
				Name:     "login",
				Usage:    "allow a user to login to the BFF via Auth0 Oauth",
				Category: "client",
				Action:   login,
				Flags:    []cli.Flag{},
			},
			{
				Name:     "status",
				Usage:    "send a status check to the BFF server",
				Category: "client",
				Action:   status,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "url",
						Aliases: []string{"u", "endpoint"},
						Usage:   "specify the URL to connect to the BFF server on",
						EnvVars: []string{"GDS_BFF_CLIENT_URL"},
						Value:   "https://bff.vaspdirectory.net",
					},
					&cli.BoolFlag{
						Name:    "nogds",
						Aliases: []string{"no-gds", "G"},
						Usage:   "health check the BFF without requesting GDS status",
						Value:   false,
					},
				},
			},
			{
				Name:     "announce",
				Usage:    "create and post an announcement",
				Category: "client",
				Action:   announce,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "url",
						Aliases: []string{"u", "endpoint"},
						Usage:   "specify the URL to connect to the BFF server on",
						EnvVars: []string{"GDS_BFF_CLIENT_URL"},
						Value:   "https://bff.vaspdirectory.net",
					},
					&cli.StringFlag{
						Name:     "token-cache",
						Aliases:  []string{"token", "t"},
						Usage:    "specify the path on disk where your access token is stored",
						EnvVars:  []string{"AUTH0_TOKEN_CACHE"},
						Required: true,
					},
				},
			},
			{
				Name:     "migrate-users",
				Usage:    "migrate Auth0 users to their organization's database record",
				Category: "client",
				Action:   migrateUsers,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "client-id",
						Aliases:  []string{"c", "id"},
						Usage:    "specify the Auth0 client ID to use for the migration",
						EnvVars:  []string{"GDS_BFF_AUTH0_CLIENT_ID"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "client-secret",
						Aliases:  []string{"s", "secret"},
						Usage:    "specify the Auth0 client secret to use for the migration",
						EnvVars:  []string{"GDS_BFF_AUTH0_CLIENT_SECRET"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "domain",
						Aliases:  []string{"d"},
						Usage:    "specify the Auth0 domain to use for the migration",
						EnvVars:  []string{"GDS_BFF_AUTH0_DOMAIN"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "endpoint",
						Aliases:  []string{"u", "url"},
						Usage:    "specify the URL to the trtl server",
						EnvVars:  []string{"GDS_BFF_DATABASE_URL"},
						Required: true,
					},
					&cli.StringFlag{
						Name:    "insecure",
						Aliases: []string{"S"},
						Usage:   "specify whether to skip TLS verification when connecting to the database",
						EnvVars: []string{"GDS_BFF_DATABASE_INSECURE"},
					},
					&cli.StringFlag{
						Name:    "cert-path",
						Aliases: []string{"C"},
						Usage:   "specify the path to the certs to use when connecting to the database",
						EnvVars: []string{"GDS_BFF_DATABASE_CERT_PATH"},
					},
					&cli.StringFlag{
						Name:    "pool-path",
						Aliases: []string{"p"},
						Usage:   "specify the path to the certs pool to use when connecting to the database",
						EnvVars: []string{"GDS_BFF_DATABASE_POOL_PATH"},
					},
					&cli.BoolFlag{
						Name:    "dry-run",
						Aliases: []string{"D", "dryrun"},
						Usage:   "specify whether to run the migration in dry-run mode, which will not write any changes to the database",
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

// Serve the GDS BFF service
func serve(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if addr := c.String("addr"); addr != "" {
		conf.BindAddr = addr
	}

	var srv *bff.Server
	if srv, err = bff.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

// Validate checks the current BFF configuration and prints the status.
func validate(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	return printJSON(conf)
}

// Login fetches an auth0 token using three-legged oauth
func login(c *cli.Context) (err error) {
	// Create a new clive server to handle the auth0 callback
	var conf clive.Config
	if conf, err = clive.NewConfig(); err != nil {
		return cli.Exit(err, 1)
	}

	var srv *clive.Server
	if srv, err = clive.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	// Get URL to redirect the user to
	var link *url.URL
	if link, err = srv.GetAuthenticationURL(); err != nil {
		return cli.Exit(err, 1)
	}

	// Open the browser window to the link
	openBrowser(link)
	fmt.Printf("To complete authentication you'll need to login with Auth0.\nIf a browser window is not automatically opened, please copy and paste the following\nlink into your browser:\n\n%s\n\n", link)

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

// Status checks if the GDS BFF is up
func status(c *cli.Context) (err error) {
	var client api.BFFClient
	if client, err = api.New(c.String("url")); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var params *api.StatusParams
	if c.Bool("nogds") {
		params = &api.StatusParams{
			NoGDS: true,
		}
	}

	var rep *api.StatusReply
	if rep, err = client.Status(ctx, params); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

// Announce creates a network announcement and posts it to the BFF.
func announce(c *cli.Context) (err error) {
	// Create the credentials to authenticate to the server.
	creds := &api.LocalCredentials{Path: c.String("token-cache")}
	if err = creds.Load(); err != nil {
		return cli.Exit(fmt.Errorf("could not load access token (run login first): %s", err), 1)
	}

	// Read the announcement from stdin
	announcement := &models.Announcement{}
	announcement.Title = readInput("Enter title: ", false)
	if len(announcement.Title) == 0 {
		return cli.Exit("please supply a post title", 1)
	}

	announcement.Body = readInput("\nPlease enter your announcement (double enter to submit, CTRL+C to quit):\n\n", true)
	if len(announcement.Body) == 0 {
		return cli.Exit("please supply an announcement to post", 1)
	}

	var client api.BFFClient
	if client, err = api.New(c.String("url"), api.WithCredentials(creds)); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = client.Login(ctx, nil); err != nil {
		return cli.Exit(err, 1)
	}

	if err = client.MakeAnnouncement(ctx, announcement); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println("announcement successfully posted!")
	return nil
}

// Migrate Auth0 users to their organization's database record.
func migrateUsers(c *cli.Context) (err error) {
	// Don't write to the database if the dry-run flag is set
	dryRun := c.Bool("dry-run")

	// Create a new Auth0 client
	authConf := config.AuthConfig{
		ClientID:     c.String("client-id"),
		ClientSecret: c.String("client-secret"),
		Domain:       c.String("domain"),
	}
	var auth0 *management.Management
	if auth0, err = auth.NewManagementClient(authConf); err != nil {
		return cli.Exit(err, 1)
	}

	// Create a new database client
	dbConf := storeconfig.StoreConfig{
		URL:      c.String("url"),
		Insecure: c.Bool("insecure"),
		CertPath: c.String("cert-path"),
		PoolPath: c.String("pool-path"),
	}
	var db store.Store
	if db, err = store.Open(dbConf); err != nil {
		return cli.Exit(err, 1)
	}

	// List Auth0 users
	var users *management.UserList
	if users, err = auth0.User.List(); err != nil {
		return cli.Exit(err, 1)
	}
	for _, user := range users.Users {
		// Get the user's organization
		appdata := &auth.AppMetadata{}
		if err = appdata.Load(user.AppMetadata); err != nil {
			return cli.Exit(err, 1)
		}

		if appdata.OrgID != "" {
			// Retrieve user's organization from the database
			var orgID uuid.UUID
			if orgID, err = models.ParseOrgID(appdata.OrgID); err != nil {
				return cli.Exit(err, 1)
			}
			var org *models.Organization
			if org, err = db.RetrieveOrganization(orgID); err != nil {
				if errors.Is(err, storeerrors.ErrEntityNotFound) {
					fmt.Printf("found user %s with missing organization %s\n", *user.Email, appdata.OrgID)
					org = &models.Organization{
						Id: appdata.OrgID,
					}

					if !dryRun {
						if _, err = db.CreateOrganization(org); err != nil {
							return cli.Exit(err, 1)
						}
						fmt.Printf("created organization %s in database for user %s\n", appdata.OrgID, *user.Email)
					}
				} else {
					return cli.Exit(err, 1)
				}
			}

			// Update the user's organization in the database
			if org.GetCollaborator(*user.Email) == nil {
				fmt.Printf("user %s is not a collaborator in their organization: %s\n", *user.Email, org.Id)
				collab := &models.Collaborator{
					Email:  *user.Email,
					UserId: *user.ID,
				}

				if err = org.AddCollaborator(collab); err != nil {
					return cli.Exit(err, 1)
				}
				fmt.Printf("created user %s as collaborator in organization %s with collab id %s\n", *user.Email, org.Id, collab.Id)

				if !dryRun {
					if err = db.UpdateOrganization(org); err != nil {
						return cli.Exit(err, 1)
					}
					fmt.Printf("updated organization %s in the database\n", org.Id)
				}
				fmt.Println()
			}
		}
	}
	return nil
}

//===========================================================================
// Helper Functions
//===========================================================================

func printJSON(msg interface{}) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(msg, "", "  "); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(data))
	return nil
}

func openBrowser(link *url.URL) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", link.String()).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", link.String()).Start()
	case "darwin":
		err = exec.Command("open", link.String()).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func readInput(prompt string, multiline bool) string {
	arr := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)

	for {
		scanner.Scan()
		text := strings.TrimSpace(scanner.Text())
		if len(text) != 0 {
			arr = append(arr, text)
			if !multiline {
				break
			}
		} else {
			break
		}
	}
	return strings.Join(arr, "\n")
}
