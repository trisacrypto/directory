package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/auth0/go-auth0/management"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var (
	db    store.Store
	auth0 *management.Management
	conf  config.Config
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
	// TODO: do we need to connect to the mainnet and testnet databases?
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

//===========================================================================
// CLI Commands
//===========================================================================

func listOrgs(c *cli.Context) (err error) {
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

	if err = db.UpdateOrganization(org); err != nil {
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
				if err = db.DeleteOrganization(curOrg.UUID()); err != nil {
					return cli.Exit(fmt.Errorf("could not delete user's current organization: %w", err), 1)
				}
			} else {
				if err = db.UpdateOrganization(curOrg); err != nil {
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

	var org *models.Organization
	if org, err = db.RetrieveOrganization(orgID); err != nil {
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
