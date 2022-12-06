package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var (
	db   store.Store
	conf config.Config
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "bffutil"
	app.Version = pkg.Version()
	app.Usage = "backend utilities for managing the BFF service and databases"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:   "orgs",
			Usage:  "list and view organization report",
			Action: organizations,
			Before: connectDB,
			After:  closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Aliases: []string{"i"},
					Usage:   "specify an organization to get more detailed information for",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func connectDB(c *cli.Context) (err error) {
	// suppress zerolog output from the store
	logger.Discard()

	// Load the configuration from the environment
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	conf.Database.ReindexOnBoot = false
	conf.ConsoleLog = false

	// Connect to the BFF main database
	// TODO: do we need to connect to the mainnet and testnet databases?
	if db, err = store.Open(conf.Database); err != nil {
		if serr, ok := status.FromError(err); ok {
			return cli.Exit(fmt.Errorf("could not open store: %s", serr.Message()), 1)
		}
		return cli.Exit(err, 1)
	}

	return nil
}

func closeDB(c *cli.Context) (err error) {
	if err = db.Close(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func organizations(c *cli.Context) (err error) {
	if id := c.String("id"); id != "" {
		// Print organization detail rather than list organizations
		var orgID uuid.UUID
		if orgID, err = models.ParseOrgID(id); err != nil {
			return cli.Exit(err, 1)
		}

		var org *models.Organization
		if org, err = db.RetrieveOrganization(orgID); err != nil {
			return cli.Exit(err, 1)
		}
		return printJSON(org)
	}

	// Create organizations report
	orgs := make([]map[string]interface{}, 0)
	iter := db.ListOrganizations()
	defer iter.Release()
	for iter.Next() {
		var org *models.Organization
		if org, err = iter.Organization(); err != nil {
			return cli.Exit(err, 1)
		}

		item := make(map[string]interface{})
		item["id"] = org.Id
		item["name"] = org.ResolveName()
		item["domain"] = org.Domain
		item["created"] = org.Created
		item["modified"] = org.Modified
		item["nCollaborators"] = len(org.Collaborators)

		orgs = append(orgs, item)
	}

	if err = iter.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(orgs)
}

//===========================================================================
// Helper Functions
//===========================================================================

func printJSON(msg interface{}) (err error) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(msg); err != nil {
		return cli.Exit(err, 1)
	}
	fmt.Println("")
	return nil
}
