package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/auth0/go-auth0/management"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var (
	db        store.Store
	mainnetDB store.Store
	testnetDB store.Store
	auth0     *management.Management
	conf      config.Config
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
			Name:   "orgs:list",
			Usage:  "list a summary of all organizations in the bff database",
			Action: listOrgs,
			Before: connectDB,
			After:  closeDB,
			Flags:  []cli.Flag{},
		},
		{
			Name:      "orgs:detail",
			Usage:     "list a summary of all organizations in the bff database",
			Action:    detailOrgs,
			ArgsUsage: "orgID [orgID ...]",
			Before:    connectDB,
			After:     closeDB,
			Flags:     []cli.Flag{},
		},
		{
			Name:   "orgs:missing",
			Usage:  "list GDS registrations that are missing organizations",
			Action: missingOrgs,
			Before: Before(loadConf, connectDB, connectGDSDatabases),
			After:  After(closeDB, closeGDSDatabases),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "no-testnet",
					Aliases: []string{"T"},
					Usage:   "don't lookup TestNet registrations in report",
				},
				&cli.BoolFlag{
					Name:    "no-mainnet",
					Aliases: []string{"M"},
					Usage:   "don't lookup MainNet registrations in report",
				},
			},
		},
		{
			Name:   "orgs:create",
			Usage:  "create an organization from existing GDS records",
			Action: createOrgs,
			Before: Before(loadConf, connectDB, connectGDSDatabases, connectAuth0),
			After:  After(closeDB, closeGDSDatabases),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "the name of the organization to create",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "domain",
					Aliases:  []string{"d"},
					Usage:    "the domain name of the organization",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "testnet-id",
					Aliases: []string{"t"},
					Usage:   "the VASP ID of the TestNet record",
				},
				&cli.StringFlag{
					Name:    "mainnet-id",
					Aliases: []string{"m"},
					Usage:   "the VASP ID of the MainNet record",
				},
				&cli.StringFlag{
					Name:     "user",
					Aliases:  []string{"u"},
					Usage:    "the auth0 user ID to add as organization leader",
					Required: true,
				},
			},
		},
		{
			Name:      "orgs:update",
			Usage:     "update an organizations name and domain",
			Action:    updateOrgs,
			ArgsUsage: "orgID",
			Before:    connectDB,
			After:     closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "update the name of the organization",
				},
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "update the domain of the organization",
				},
			},
		},
		{
			Name:      "orgs:rmsub",
			Usage:     "remove an organization's registration record for either testnet or mainnet",
			Action:    rmsubOrgs,
			ArgsUsage: "orgID",
			Before:    connectDB,
			After:     closeDB,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "mainnet",
					Aliases: []string{"m"},
					Usage:   "delete mainnet registration record",
				},
				&cli.BoolFlag{
					Name:    "testnet",
					Aliases: []string{"t"},
					Usage:   "delete testnet registration record",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "do not prompt to confirm operation",
				},
			},
		},
		{
			Name:      "orgs:delete",
			Usage:     "delete an organization and remove it from its users",
			Action:    deleteOrgs,
			ArgsUsage: "orgID",
			Before:    Before(loadConf, connectDB, connectAuth0),
			After:     closeDB,
			Flags:     []cli.Flag{},
		},
		{
			Name:   "orgs:cleanup",
			Usage:  "removes any organizations that have zero collaborators",
			Action: cleanupOrgs,
			Before: connectDB,
			After:  closeDB,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "do not prompt to confirm org delete",
				},
			},
		},
		{
			Name:   "collabs:add",
			Usage:  "add a collaborator to an organization",
			Action: addCollab,
			Before: Before(loadConf, connectDB, connectAuth0),
			After:  closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "org",
					Aliases:  []string{"o"},
					Usage:    "specify the organization id to add the collaborator to",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "user",
					Aliases:  []string{"u"},
					Usage:    "specify the auth0 id of the user to make a collaborator",
					Required: true,
				},
			},
		},
		{
			Name:   "collabs:delete",
			Usage:  "remove a collaborator from an organization",
			Action: deleteCollab,
			Before: Before(loadConf, connectDB, connectAuth0),
			After:  closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "org",
					Aliases:  []string{"o"},
					Usage:    "specify the organization id to add the collaborator to",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "user",
					Aliases:  []string{"u"},
					Usage:    "specify the auth0 id of the user to make a collaborator",
					Required: true,
				},
			},
		},
		{
			Name:   "appdata:sortorgs",
			Usage:  "sort the organization list on each user's app_metadata",
			Action: sortAppdataOrgs,
			Before: Before(loadConf, connectAuth0),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "display user app_metadata changes without updating them",
				},
			},
		},
		{
			Name:   "appdata:dedupeorgs",
			Usage:  "remove duplicate organizations from a user's app_metadata",
			Action: dedupeAppdataOrgs,
			Before: Before(loadConf, connectAuth0),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dry-run",
					Aliases: []string{"d"},
					Usage:   "display user duplicate orgs without removing them",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

//===========================================================================
// Before/After CLI Commands
//===========================================================================

func loadConf(c *cli.Context) (err error) {
	// suppress zerolog output from the store
	logger.Discard()

	// Load the configuration from the environment
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	conf.Database.ReindexOnBoot = false
	conf.ConsoleLog = false
	return nil
}

func connectDB(c *cli.Context) (err error) {
	if conf.IsZero() {
		if err = loadConf(c); err != nil {
			return err
		}
	}

	// Connect to the BFF main database
	if db, err = store.Open(conf.Database); err != nil {
		if serr, ok := status.FromError(err); ok {
			return cli.Exit(fmt.Errorf("could not open store: %s", serr.Message()), 1)
		}
		return cli.Exit(err, 1)
	}

	return nil
}

func connectAuth0(c *cli.Context) (err error) {
	if conf.IsZero() {
		if err = loadConf(c); err != nil {
			return err
		}
	}

	if auth0, err = auth.NewManagementClient(conf.Auth0); err != nil {
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

func connectGDSDatabases(c *cli.Context) (err error) {
	if conf.IsZero() {
		if err = loadConf(c); err != nil {
			return err
		}
	}

	// Connect to the GDS TestNet database
	if testnetDB, err = store.Open(conf.TestNet.Database); err != nil {
		if serr, ok := status.FromError(err); ok {
			return cli.Exit(fmt.Errorf("could not open testnet store: %s", serr.Message()), 1)
		}
		return cli.Exit(err, 1)
	}

	// Connect to the GDS MainNet database
	if mainnetDB, err = store.Open(conf.MainNet.Database); err != nil {
		if serr, ok := status.FromError(err); ok {
			return cli.Exit(fmt.Errorf("could not open mainnet store: %s", serr.Message()), 1)
		}
		return cli.Exit(err, 1)
	}

	return nil
}

func closeGDSDatabases(c *cli.Context) (err error) {
	if mainnetDB != nil {
		if dberr := mainnetDB.Close(); dberr != nil {
			err = multierror.Append(err, dberr)
		}
	}

	if testnetDB != nil {
		if dberr := testnetDB.Close(); dberr != nil {
			err = multierror.Append(err, dberr)
		}
	}

	return err
}

//===========================================================================
// CLI Commands
//===========================================================================

func listOrgs(c *cli.Context) (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Create organizations report
	orgs := make([]map[string]interface{}, 0)
	iter := db.ListOrganizations(ctx)
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
		item["testnet"] = org.Testnet
		item["mainnet"] = org.Mainnet

		orgs = append(orgs, item)
	}

	if err = iter.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(orgs)
}

func detailOrgs(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.Exit("specify at least one organization ID", 1)
	}

	orgs := make([]*models.Organization, 0, c.NArg())
	for i := 0; i < c.NArg(); i++ {
		var org *models.Organization
		if org, err = GetOrg(c.Args().Get(i)); err != nil {
			return cli.Exit(err, 1)
		}
		orgs = append(orgs, org)
	}

	if len(orgs) == 1 {
		return printJSON(orgs[0])
	}
	return printJSON(orgs)
}

func missingOrgs(c *cli.Context) (err error) {
	if c.Bool("no-testnet") && c.Bool("no-mainnet") {
		return cli.Exit("no GDS networks specified to analyze", 0)
	}

	// Step one: compile all directory records from existing organizations
	testnet := make(map[string]struct{})
	mainnet := make(map[string]struct{})

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	orgs := db.ListOrganizations(ctx)
	defer orgs.Release()
	for orgs.Next() {
		var org *models.Organization
		if org, err = orgs.Organization(); err != nil {
			return cli.Exit(err, 1)
		}

		if org.Testnet != nil && org.Testnet.Id != "" {
			testnet[org.Testnet.Id] = struct{}{}
		}

		if org.Mainnet != nil && org.Mainnet.Id != "" {
			mainnet[org.Mainnet.Id] = struct{}{}
		}
	}

	if err = orgs.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	// Step one and a half: create a CSV document to write records to
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"id", "name", "common name", "registered directory"})

	// Step two: loop through TestNet to see what registrations are missing
	if !c.Bool("no-testnet") {
		vasps := testnetDB.ListVASPs(ctx)
		defer vasps.Release()
		for vasps.Next() {
			var vasp *pb.VASP
			if vasp, err = vasps.VASP(); err != nil {
				return cli.Exit(err, 1)
			}

			if _, ok := testnet[vasp.Id]; !ok {
				name, _ := vasp.Name()
				row := []string{vasp.Id, name, vasp.CommonName, vasp.RegisteredDirectory}
				writer.Write(row)
			}
		}

		if err = vasps.Error(); err != nil {
			return cli.Exit(err, 1)
		}
	}

	// Step three: loop through MainNet to see what registrations are missing
	if !c.Bool("no-mainnet") {
		vasps := mainnetDB.ListVASPs(ctx)
		defer vasps.Release()
		for vasps.Next() {
			var vasp *pb.VASP
			if vasp, err = vasps.VASP(); err != nil {
				return cli.Exit(err, 1)
			}

			if _, ok := mainnet[vasp.Id]; !ok {
				name, _ := vasp.Name()
				row := []string{vasp.Id, name, vasp.CommonName, vasp.RegisteredDirectory}
				writer.Write(row)
			}
		}

		if err = vasps.Error(); err != nil {
			return cli.Exit(err, 1)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func createOrgs(c *cli.Context) (err error) {
	var (
		mainnetVASP *pb.VASP
		testnetVASP *pb.VASP
		user        *management.User
		permissions *management.PermissionList
		appdata     *auth.AppMetadata
		username    string
		mainname    string
		testname    string
	)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if vaspID := c.String("mainnet-id"); vaspID != "" {
		if mainnetVASP, err = mainnetDB.RetrieveVASP(ctx, vaspID); err != nil {
			return cli.Exit(err, 1)
		}
		mainname, _ = mainnetVASP.Name()
		if mainname == "" {
			mainname = mainnetVASP.CommonName
		}
	} else {
		mainname = "N/A"
	}

	if vaspID := c.String("testnet-id"); vaspID != "" {
		if testnetVASP, err = testnetDB.RetrieveVASP(ctx, vaspID); err != nil {
			return cli.Exit(err, 1)
		}
		testname, _ = testnetVASP.Name()
		if testname == "" {
			testname = testnetVASP.CommonName
		}
	} else {
		testname = "N/A"
	}

	if user, err = auth0.User.Read(c.String("user")); err != nil {
		return cli.Exit(err, 1)
	}

	// Fetch the appdata and permissions of the user
	appdata = &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		return cli.Exit(err, 1)
	}

	if permissions, err = auth0.User.Permissions(*user.ID); err != nil {
		return cli.Exit(err, 1)
	}

	if !HasPermission(auth.SwitchOrganizations, permissions) {
		return cli.Exit("the user must be a TSP user", 1)
	}

	// Ask if we should proceed
	username, _ = auth.UserDisplayName(user)
	if !askForConfirmation(fmt.Sprintf("create org for TestNet: %s and MainNet: %s records with user %s?", testname, mainname, username)) {
		return cli.Exit("canceled at request of user", 0)
	}

	// Create new organization record
	org := &models.Organization{
		Name:      c.String("name"),
		Domain:    c.String("domain"),
		CreatedBy: "support@rotational.io",
	}

	if org.Domain, err = bff.NormalizeDomain(org.Domain); err != nil {
		return cli.Exit(err, 1)
	}

	if err = bff.ValidateDomain(org.Domain); err != nil {
		return cli.Exit(err, 1)
	}

	// TODO: Check for duplicate domains

	// Add the user as a collaborator to the organization
	// NOTE: expecting the user to be a TSP so no roles are modified
	collaborator := &models.Collaborator{
		Email:    *user.Email,
		UserId:   *user.ID,
		Verified: *user.EmailVerified,
	}

	if err = org.AddCollaborator(collaborator); err != nil {
		return cli.Exit(err, 1)
	}

	// Add the directory records
	if mainnetVASP != nil {
		org.Mainnet = &models.DirectoryRecord{
			Id:                  mainnetVASP.Id,
			RegisteredDirectory: mainnetVASP.RegisteredDirectory,
			CommonName:          mainnetVASP.CommonName,
			Submitted:           mainnetVASP.FirstListed,
		}
	}

	if testnetVASP != nil {
		org.Testnet = &models.DirectoryRecord{
			Id:                  testnetVASP.Id,
			RegisteredDirectory: testnetVASP.RegisteredDirectory,
			CommonName:          testnetVASP.CommonName,
			Submitted:           testnetVASP.FirstListed,
		}
	}

	// Create the registration form
	var vasp *pb.VASP
	if mainnetVASP != nil {
		vasp = mainnetVASP
	} else if testnetVASP != nil {
		vasp = testnetVASP
	}

	if vasp != nil {
		reg := models.NewRegisterForm()
		reg.Website = vasp.Website
		reg.BusinessCategory = vasp.BusinessCategory
		reg.VaspCategories = vasp.VaspCategories
		reg.EstablishedOn = vasp.EstablishedOn
		reg.OrganizationName = org.Name
		reg.Entity = vasp.Entity
		reg.Contacts = vasp.Contacts
		reg.Trixo = vasp.Trixo

		if mainnetVASP != nil {
			reg.Mainnet = &models.NetworkDetails{
				CommonName: mainnetVASP.CommonName,
				Endpoint:   mainnetVASP.TrisaEndpoint,
			}
		}

		if testnetVASP != nil {
			reg.Testnet = &models.NetworkDetails{
				CommonName: testnetVASP.CommonName,
				Endpoint:   testnetVASP.TrisaEndpoint,
			}
		}

		reg.State.Current = 6
		reg.State.ReadyToSubmit = reg.ReadyToSubmit("all")
		reg.State.Started = vasp.FirstListed
		reg.State.Steps = []*models.FormStep{
			{Key: 1, Status: "done"},
			{Key: 2, Status: "done"},
			{Key: 3, Status: "done"},
			{Key: 4, Status: "done"},
			{Key: 5, Status: "done"},
			{Key: 6, Status: "done"},
		}

		org.Registration = reg
	} else {
		org.Registration = models.NewRegisterForm()
	}

	// Create the organization
	if _, err = db.CreateOrganization(ctx, org); err != nil {
		return cli.Exit(err, 1)
	}

	// Add the user to the organization
	appdata.AddOrganization(org.Id)
	if err = SaveAppMetadata(*user.ID, *appdata); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func updateOrgs(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		if c.NArg() == 0 {
			return cli.Exit("specify an orgID to update", 1)
		}
		return cli.Exit("can only update one organization at a time", 1)
	}

	if c.String("name") == "" && c.String("domain") == "" {
		return cli.Exit("specify name or domain to update", 1)
	}

	var org *models.Organization
	if org, err = GetOrg(c.Args().Get(0)); err != nil {
		return cli.Exit(err, 1)
	}

	save := false
	if name := c.String("name"); name != "" && name != org.Name {
		org.Name = name
		save = true
	}

	if domain := c.String("domain"); domain != "" && domain != org.Domain {
		org.Domain = domain
		save = true
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if save {
		if err = db.UpdateOrganization(ctx, org); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		fmt.Println("no changes made to organization")
	}
	return nil
}

func rmsubOrgs(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		if c.NArg() == 0 {
			return cli.Exit("specify an orgID to modify", 1)
		}
		return cli.Exit("can only modify one organization at a time", 1)
	}

	mainnet := c.Bool("mainnet")
	testnet := c.Bool("testnet")

	if !mainnet && !testnet {
		return cli.Exit("specify either mainnet, testnet, or both to remove submission for", 1)
	}

	var org *models.Organization
	if org, err = GetOrg(c.Args().Get(0)); err != nil {
		return cli.Exit(err, 1)
	}

	save := false

	if mainnet {
		if org.Mainnet != nil {
			// Prompt for confirmation to delete
			if !c.Bool("force") && !askForConfirmation(fmt.Sprintf("delete mainnet registration from %s?", org.ResolveName())) {
				return cli.Exit("operation cancelled by user", 0)
			}
			org.Mainnet = nil
			save = true
		} else {
			fmt.Println("organization has no mainnet submission")
		}
	}

	if testnet {
		if org.Testnet != nil {
			// Prompt for confirmation to delete
			if !c.Bool("force") && !askForConfirmation(fmt.Sprintf("delete testnet registration from %s?", org.ResolveName())) {
				return cli.Exit("operation cancelled by user", 0)
			}
			org.Testnet = nil
			save = true
		} else {
			fmt.Println("organization has no testnet submission")
		}
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if save {
		if err = db.UpdateOrganization(ctx, org); err != nil {
			return cli.Exit(err, 1)
		}
	}
	return nil
}

func deleteOrgs(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		if c.NArg() == 0 {
			return cli.Exit("specify an orgID to delete", 1)
		}
		return cli.Exit("can only delete one organization at a time", 1)
	}

	var org *models.Organization
	if org, err = GetOrg(c.Args().Get(0)); err != nil {
		return cli.Exit(err, 1)
	}

	// Fetch collaborator auth0 appdata and verify that all collaborators are part of
	// at least one additional organization besides the one being deleted.
	appdata := make(map[string]*auth.AppMetadata)
	for _, collaborator := range org.Collaborators {
		if collaborator.UserId == "" {
			// Invited user who hasn't joined yet
			continue
		}

		var user *management.User
		if user, err = auth0.User.Read(collaborator.UserId); err != nil {
			return cli.Exit(fmt.Errorf("could not fetch user for %s", collaborator.Email), 1)
		}

		meta := &auth.AppMetadata{}
		if err = meta.Load(user.AppMetadata); err != nil {
			return cli.Exit(fmt.Errorf("could not load app metadata for %s", collaborator.Email), 1)
		}

		// Check to make sure the orgID isn't the only org the user belongs to.
		uorgs := make(map[string]struct{})
		for _, uorg := range meta.Organizations {
			if uorg != org.Id {
				uorgs[uorg] = struct{}{}
			}
		}

		if len(uorgs) == 0 {
			if !askForConfirmation(fmt.Sprintf("user %s only belongs to this organization, a new organization will be created the next time they login, continue?", collaborator.Email)) {
				return cli.Exit("operation cancelled by user", 0)
			}
		}

		fmt.Printf("%s will be removed from organization\n", collaborator.Email)
		appdata[collaborator.UserId] = meta
	}

	// Prompt for confirmation to delete
	if !askForConfirmation(fmt.Sprintf("delete %s with %d collaborators?", org.ResolveName(), len(org.Collaborators))) {
		return cli.Exit("operation cancelled by user", 0)
	}

	// Remove the organization from all of the collaborators
	for uid, umeta := range appdata {
		umeta.RemoveOrganization(org.Id)
		if umeta.OrgID == org.Id {
			if len(umeta.Organizations) > 0 {
				umeta.OrgID = umeta.Organizations[0]
			} else {
				umeta.OrgID = ""
			}
		}

		if err = SaveAppMetadata(uid, *umeta); err != nil {
			return cli.Exit(err, 1)
		}
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Last step: delete the organization from the database
	if err = db.DeleteOrganization(ctx, org.UUID()); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func cleanupOrgs(c *cli.Context) (err error) {
	orgsDeleted := 0
	force := c.Bool("force")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	iter := db.ListOrganizations(ctx)
	defer iter.Release()
	for iter.Next() {
		var org *models.Organization
		if org, err = iter.Organization(); err != nil {
			return cli.Exit(err, 1)
		}

		if len(org.Collaborators) == 0 {
			if !force && !askForConfirmation(fmt.Sprintf("org %s has 0 collaborators, delete?", org.ResolveName())) {
				continue
			}

			if err = db.DeleteOrganization(ctx, org.UUID()); err != nil {
				return cli.Exit(fmt.Errorf("could not delete %s: %w", org.ResolveName(), err), 1)
			}
			orgsDeleted++
		}
	}

	if err = iter.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("deleted %d organizations\n", orgsDeleted)
	return nil
}

func addCollab(c *cli.Context) (err error) {
	// Collect the organization from the database
	var org *models.Organization
	if org, err = GetOrg(c.String("org")); err != nil {
		return cli.Exit(err, 1)
	}

	// Fetch the user from auth0.
	var user *management.User
	if user, err = auth0.User.Read(c.String("user")); err != nil {
		return cli.Exit(err, 1)
	}

	// Get the user appdata, roles, and permissions
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		return cli.Exit(err, 1)
	}

	var roles *management.RoleList
	if roles, err = auth0.User.Roles(*user.ID); err != nil {
		return cli.Exit(err, 1)
	}

	var permissions *management.PermissionList
	if permissions, err = auth0.User.Permissions(*user.ID); err != nil {
		return cli.Exit(err, 1)
	}

	// Check if the user is already in the organization
	if appdata.OrgID == org.Id {
		return cli.Exit("user is already a collaborator in this organization", 0)
	}

	// Ask if we should proceed
	username, _ := auth.UserDisplayName(user)
	fmt.Printf("User %s (%s, switch_organizations=%t) has appdata.OrgID %q\n", username, StringifyRoles(roles), HasPermission(auth.SwitchOrganizations, permissions), appdata.OrgID)
	if !askForConfirmation(fmt.Sprintf("add user %q to organization %q?", username, org.ResolveName())) {
		return cli.Exit("canceled at request of user", 0)
	}

	// Add the user to the new organization
	collaborator := &models.Collaborator{
		Email:    *user.Email,
		UserId:   *user.ID,
		Verified: *user.EmailVerified,
	}

	if err = org.AddCollaborator(collaborator); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if err = db.UpdateOrganization(ctx, org); err != nil {
		return cli.Exit(fmt.Errorf("could not update organization: %w", err), 1)
	}

	if HasPermission(auth.SwitchOrganizations, permissions) {
		// If the user has the switch organizations permission then we just add them to
		// the new organization but do not remove them from the old organization
		appdata.AddOrganization(org.Id)
	} else {
		// If the user does not have the switch organizations permission, remove them
		// from their current organization and add them to the new organization.
		if appdata.OrgID != "" {
			var curOrg *models.Organization
			if curOrg, err = GetOrg(appdata.OrgID); err != nil {
				return cli.Exit(err, 1)
			}

			curOrg.DeleteCollaborator(*user.Email)

			if len(curOrg.Collaborators) == 0 {
				// If there are no more collaborators in the current org, delete it.
				if err = db.DeleteOrganization(ctx, curOrg.UUID()); err != nil {
					return cli.Exit(fmt.Errorf("could not delete user's current organization: %w", err), 1)
				}
			} else {
				if err = db.UpdateOrganization(ctx, curOrg); err != nil {
					return cli.Exit(fmt.Errorf("could not update user's current organization: %w", err), 1)
				}
			}
		}
	}

	// Update user's app metadata to reflect the user's currently selected organization.
	appdata.UpdateOrganization(org)
	if err = SaveAppMetadata(*user.ID, *appdata); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func deleteCollab(c *cli.Context) (err error) {
	// Collect the organization from the database
	var org *models.Organization
	if org, err = GetOrg(c.String("org")); err != nil {
		return cli.Exit(err, 1)
	}

	// Fetch the user from auth0.
	var user *management.User
	if user, err = auth0.User.Read(c.String("user")); err != nil {
		return cli.Exit(err, 1)
	}

	// Get the user appdata, roles, and permissions
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		return cli.Exit(err, 1)
	}

	var roles *management.RoleList
	if roles, err = auth0.User.Roles(*user.ID); err != nil {
		return cli.Exit(err, 1)
	}

	var permissions *management.PermissionList
	if permissions, err = auth0.User.Permissions(*user.ID); err != nil {
		return cli.Exit(err, 1)
	}

	if !HasPermission(auth.SwitchOrganizations, permissions) {
		return cli.Exit("cannot remove a collaborator without TRISA Service Provider role", 1)
	}

	if len(appdata.Organizations) < 2 {
		return cli.Exit("cannot remove collaborator without another organization to fall back on", 1)
	}

	// Ask if we should proceed
	username, _ := auth.UserDisplayName(user)
	fmt.Printf("User %s (%s, switch_organizations=%t) has appdata.OrgID %q\n", username, StringifyRoles(roles), HasPermission(auth.SwitchOrganizations, permissions), appdata.OrgID)
	if !askForConfirmation(fmt.Sprintf("add user %q to organization %q?", username, org.ResolveName())) {
		return cli.Exit("canceled at request of user", 0)
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Remove collaborator from the organization (won't error if not exists)
	org.DeleteCollaborator(*user.Email)
	if err = db.UpdateOrganization(ctx, org); err != nil {
		return cli.Exit(fmt.Errorf("could not update organization: %w", err), 1)
	}

	appdata.RemoveOrganization(org.Id)
	if appdata.OrgID == org.Id {
		appdata.OrgID = appdata.Organizations[0]
	}

	if err = SaveAppMetadata(*user.ID, *appdata); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func sortAppdataOrgs(c *cli.Context) (err error) {
	// Get all users in the tenant
	var users *management.UserList
	if users, err = auth0.User.List(); err != nil {
		return cli.Exit(err, 1)
	}

	// Recreate the org list for each user
	var updated int
	for _, user := range users.Users {
		appdata := &auth.AppMetadata{}
		if err = appdata.Load(user.AppMetadata); err != nil {
			return cli.Exit(err, 1)
		}

		// Get the user's current org list
		orgs := appdata.GetOrganizations()
		if len(orgs) == 0 {
			continue
		}

		fmt.Printf("current org list for user %s (%s): %v\n", *user.Email, *user.ID, orgs)

		// Create the sorted list, using the internal method which ensures that the
		// organizations are sorted in the intended order
		appdata.Organizations = []string{}
		for _, org := range orgs {
			appdata.AddOrganization(org)
		}
		fmt.Printf("sorted org list for user %s (%s): %v\n", *user.Email, *user.ID, appdata.GetOrganizations())

		// Update the user's app metadata on the Auth0 tenant
		if !c.Bool("dry-run") {
			if err = SaveAppMetadata(*user.ID, *appdata); err != nil {
				return cli.Exit(err, 1)
			}
			fmt.Printf("updated app metadata for user %s (%s)\n", *user.Email, *user.ID)
			updated++
		}

		fmt.Println()
	}

	fmt.Printf("updated app metadata for %d users\n", updated)

	return nil
}

func dedupeAppdataOrgs(c *cli.Context) (err error) {
	// Get all users in the tenant
	var users *management.UserList
	if users, err = auth0.User.List(); err != nil {
		return cli.Exit(err, 1)
	}

	// Recreate org list for each user
	for _, user := range users.Users {
		appdata := &auth.AppMetadata{}
		if err = appdata.Load(user.AppMetadata); err != nil {
			return cli.Exit(err, 1)
		}

		// Get the user's current org list
		orgs := appdata.GetOrganizations()
		if len(orgs) == 0 {
			continue
		}

		// Check for duplicate orgs and remove them.
		seen := make(map[string]struct{})
		appdata.Organizations = []string{}
		for _, org := range orgs {
			if _, ok := seen[org]; !ok {
				seen[org] = struct{}{}
				appdata.Organizations = append(appdata.Organizations, org)
			} else {
				fmt.Printf("found duplicate org %q from user %s (%s)\n", org, *user.Email, *user.ID)
			}
		}

		fmt.Printf("current org list for user %s (%s): %v\n", *user.Email, *user.ID, appdata.Organizations)

		if !c.Bool("dry-run") {
			if err = SaveAppMetadata(*user.ID, *appdata); err != nil {
				return cli.Exit(err, 1)
			}
		}
	}

	return nil
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

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not interpret response: %s", err)
			os.Exit(1)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func GetOrg(id string) (_ *models.Organization, err error) {
	var orgID uuid.UUID
	if orgID, err = models.ParseOrgID(id); err != nil {
		return nil, err
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	var org *models.Organization
	if org, err = db.RetrieveOrganization(ctx, orgID); err != nil {
		return nil, err
	}
	return org, nil
}

func Before(funcs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, f := range funcs {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
}

func After(funcs ...cli.AfterFunc) cli.AfterFunc {
	return func(c *cli.Context) error {
		for _, f := range funcs {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
}

func StringifyRoles(roles *management.RoleList) string {
	switch roles.Total {
	case 0:
		return "no roles"
	case 1:
		return fmt.Sprintf("role %s", roles.Roles[0].GetName())
	default:
		names := make([]string, 0, roles.Total)
		for _, role := range roles.Roles {
			names = append(names, role.GetName())
		}
		return fmt.Sprintf("roles %s", strings.Join(names, ", "))
	}
}

func HasPermission(perm string, permissions *management.PermissionList) bool {
	for _, permission := range permissions.Permissions {
		if permission.GetName() == perm {
			return true
		}
	}
	return false
}

func SaveAppMetadata(uid string, appdata auth.AppMetadata) (err error) {
	// Create a blank user with no data but the app data
	user := &management.User{}
	if user.AppMetadata, err = appdata.Dump(); err != nil {
		return err
	}

	if err = auth0.User.Update(uid, user); err != nil {
		return err
	}
	return nil
}
