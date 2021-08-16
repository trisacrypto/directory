package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/urfave/cli"
)

var (
	api     *sectigo.Sectigo
	encoder *json.Encoder
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "sectigo"
	app.Version = pkg.Version()
	app.Usage = "CLI helper for Sectigo API access and debugging"
	app.Before = initAPI
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "u, username",
			Usage:  "API access login username",
			EnvVar: sectigo.UsernameEnv,
		},
		cli.StringFlag{
			Name:   "p, password",
			Usage:  "API access login password",
			EnvVar: sectigo.PasswordEnv,
		},
		cli.StringFlag{
			Name:   "P, profile",
			Usage:  "Specify the profile to use for certificate requests",
			EnvVar: sectigo.ProfileEnv,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "auth",
			Usage:  "check authentication status with server",
			Action: auth,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "d, debug",
					Usage: "do not refresh or authenticate, print state and exit",
				},
				cli.BoolFlag{
					Name:  "C, cache",
					Usage: "print cache location and exit",
				},
			},
		},
		{
			Name:   "create",
			Usage:  "create single certificate batch",
			Action: createSingle,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "a, authority",
					Usage: "id of the authority or profile to issue the cert",
				},
				cli.StringFlag{
					Name:  "d, domain",
					Usage: "common name of the subject to issue the cert for",
				},
				cli.StringFlag{
					Name:  "p, password",
					Usage: "password for script (automatically generated by default)",
				},
				cli.StringFlag{
					Name:  "c, country",
					Usage: "two-letter ISO 3166-1 country code",
				},
				cli.StringFlag{
					Name:  "s, state",
					Usage: "state or province name",
				},
				cli.StringFlag{
					Name:  "l, locality",
					Usage: "city or locality name",
				},
				cli.StringFlag{
					Name:  "b, batch-name",
					Usage: "description of the batch for review purposes",
				},
			},
		},
		{
			Name:      "upload",
			Usage:     "upload a certificate signing request",
			ArgsUsage: "request.csr",
			Action:    uploadCSR,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "p, profile",
					Usage: "profile/authority id of the cert issuer",
				},
			},
		},
		{
			Name:   "batches",
			Usage:  "view batch jobs for certificate creation",
			Action: batches,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "i, id",
					Usage: "if specified get detail for batch with id",
				},
				cli.BoolFlag{
					Name:  "s, status",
					Usage: "get batch processing status",
				},
			},
		},
		{
			Name:   "download",
			Usage:  "download batch as a ZIP file",
			Action: download,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "i, id",
					Usage: "the batch ID to download",
				},
				cli.StringFlag{
					Name:  "o, outdir",
					Usage: "the directory to download the zip file to",
				},
			},
		},
		{
			Name:   "licenses",
			Usage:  "view the ordered/issued certificates",
			Action: licenses,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "authorities",
			Usage:  "view the current users authorities by ecosystem",
			Action: authorities,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "b, balances",
					Usage: "also fetch balance for each authority",
				},
			},
		},
		{
			Name:   "profiles",
			Usage:  "view profiles available to the user",
			Action: profiles,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "i, id",
					Usage: "if specified get detail for profile with id",
				},
				cli.BoolFlag{
					Name:  "p, params",
					Usage: "if specified, get params for profile with id",
				},
			},
		},
		{
			Name:   "organization",
			Usage:  "view organization of the current user",
			Action: organization,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "find",
			Usage:  "search for certs by common name and serial number",
			Action: findCert,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "n, common-name",
					Usage: "search by common name",
				},
				cli.StringFlag{
					Name:  "s, serial-number",
					Usage: "search by serial number",
				},
			},
		},
		{
			Name:   "revoke",
			Usage:  "revoke a certificate by serial number",
			Action: revokeCert,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "p, profile",
					Usage: "profile/authority id of the cert issuer",
				},
				cli.StringFlag{
					Name:  "r, reason",
					Usage: "RFC 5280 reason text",
				},
				cli.StringFlag{
					Name:  "s, serial-number",
					Usage: "serial number of the cert to revoke",
				},
			},
		},
	}

	app.Run(os.Args)
}

func initAPI(c *cli.Context) (err error) {
	if api, err = sectigo.New(c.String("username"), c.String("password"), c.String("profile")); err != nil {
		return cli.NewExitError(err, 1)
	}

	encoder = json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return nil
}

func auth(c *cli.Context) (err error) {
	creds := api.Creds()

	if c.Bool("cache") {
		if cacheFile := creds.CacheFile(); cacheFile != "" {
			fmt.Println(cacheFile)
		} else {
			fmt.Println("no credentials cache file exists")
		}
		return nil
	}

	if c.Bool("debug") {
		if creds.Valid() {
			fmt.Printf("credentials are valid until %s\n", creds.ExpiresAt)
			return nil
		}

		if creds.Current() {
			fmt.Printf("credentials are current until %s\n", creds.RefreshBy)
			return nil
		}

		fmt.Println("credentials are expired or invalid")
		return nil
	}

	if !creds.Valid() {
		if creds.Refreshable() {
			if err = api.Refresh(); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			if err = api.Authenticate(); err != nil {
				return cli.NewExitError(err, 1)
			}
		}
	}

	fmt.Println("user authenticated and credentials cached")
	return nil
}

func createSingle(c *cli.Context) (err error) {
	domain := c.String("domain")
	if domain == "" {
		return cli.NewExitError("must specify domain name of cert subject", 1)
	}

	authority := c.Int("authority")
	if authority == 0 {
		return cli.NewExitError("must specify authority ID", 1)
	}

	params := make(map[string]string)
	params["commonName"] = domain
	params["dNSName"] = domain
	params["pkcs12Password"] = c.String("password")

	if params["pkcs12Password"] == "" {
		params["pkcs12Password"] = randomPassword(10)
		fmt.Printf("pkcs12 password: %s\n", params["pkcs12Password"])
	}

	if api.Profile() == sectigo.ProfileCipherTraceEndEntityCertificate || api.Profile() == sectigo.ProfileIDCipherTraceEndEntityCertificate {
		params["localityName"] = c.String("locality")
		params["stateOrProvinceName"] = c.String("state")
		params["countryName"] = c.String("country")
	}

	batchName := c.String("batch-name")
	if batchName == "" {
		batchName = fmt.Sprintf("new certs for %s", domain)
	}

	var rep *sectigo.BatchResponse
	if rep, err = api.CreateSingleCertBatch(authority, batchName, params); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)
	return nil
}

func uploadCSR(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		return cli.NewExitError("specify the path to one CSR for upload", 1)
	}

	path := c.Args().First()
	filename := filepath.Base(path)

	// Create the request parameters
	profileId := c.Int("profile")
	params := make(map[string]string)

	if profileId == 0 {
		return cli.NewExitError("specify the profile ID to sign the cert under", 1)
	}

	// Load the CSR data from the file
	var csrData []byte
	if csrData, err = ioutil.ReadFile(path); err != nil {
		return cli.NewExitError(fmt.Errorf("could not read %s: %s", path, err), 1)
	}

	var rep *sectigo.BatchResponse
	if rep, err = api.UploadCSRBatch(profileId, filename, csrData, params); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)
	return nil
}

func batches(c *cli.Context) (err error) {
	id := c.Int("id")
	if id != 0 {
		// Perform batch detail lookup
		if c.Bool("status") {
			var rep *sectigo.ProcessingInfoResponse
			if rep, err = api.ProcessingInfo(id); err != nil {
				return cli.NewExitError(err, 1)
			}

			printJSON(rep)
			return nil
		}

		var rep *sectigo.BatchResponse
		if rep, err = api.BatchDetail(id); err != nil {
			return cli.NewExitError(err, 1)
		}

		printJSON(rep)
		return nil
	}

	return cli.NewExitError("specify batch id to get information", 1)
}

func download(c *cli.Context) (err error) {
	outdir := c.String("outdir")
	batch := c.Int("id")
	if batch == 0 {
		return cli.NewExitError("must specify batch id for download", 1)
	}

	var path string
	if path, err = api.Download(batch, outdir); err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("downloaded batch %d to %s\n", batch, path)
	fmt.Println("after unzipping, unencrypt with your password using `openssl pkcs12 -in INFILE.p12 -out OUTFILE.crt -nodes`")
	return nil
}

func licenses(c *cli.Context) (err error) {
	var rep *sectigo.LicensesUsedResponse
	if rep, err = api.LicensesUsed(); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)
	return nil
}

func authorities(c *cli.Context) (err error) {
	var rep []*sectigo.AuthorityResponse
	if rep, err = api.UserAuthorities(); err != nil {
		return cli.NewExitError(err, 1)
	}

	// Print the authority details
	printJSON(rep)

	if c.Bool("balances") {
		// Fetch the balances for each authority and print them
		balances := make(map[int]int)
		for _, authority := range rep {
			if balances[authority.ID], err = api.AuthorityAvailableBalance(authority.ID); err != nil {
				return cli.NewExitError(err, 1)
			}
		}
		printJSON(balances)
	}

	return nil
}

func profiles(c *cli.Context) (err error) {
	pid := c.Int("id")
	if pid != 0 {
		// Perform a detail request instead of a list request
		// Get params detail for the profile
		if c.Bool("params") {
			var rep []*sectigo.ProfileParamsResponse
			if rep, err = api.ProfileParams(pid); err != nil {
				return cli.NewExitError(err, 1)
			}
			printJSON(rep)
			return nil
		}

		// Get profile extended information
		var rep *sectigo.ProfileDetailResponse
		if rep, err = api.ProfileDetail(pid); err != nil {
			return cli.NewExitError(err, 1)
		}
		printJSON(rep)
		return nil
	}

	if c.Bool("params") {
		return cli.NewExitError("must specify id to get profile params", 1)
	}

	var rep []*sectigo.ProfileResponse
	if rep, err = api.Profiles(); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)
	return nil
}

func organization(c *cli.Context) (err error) {
	var rep *sectigo.OrganizationResponse
	if rep, err = api.Organization(); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)
	return nil
}

func findCert(c *cli.Context) (err error) {
	var rep *sectigo.FindCertificateResponse
	if rep, err = api.FindCertificate(c.String("common-name"), c.String("serial-number")); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)
	return nil
}

func revokeCert(c *cli.Context) (err error) {
	pid := c.Int("profile")
	if pid == 0 {
		return cli.NewExitError("must specify profile id", 1)
	}

	var reasonCode sectigo.CRLReason
	if reasonCode, err = sectigo.RevokeReasonCode(c.String("reason")); err != nil {
		return cli.NewExitError(err, 1)
	}

	sn := c.String("serial-number")
	if sn == "" {
		return cli.NewExitError("must specify serial number of certificate", 1)
	}

	if err = api.RevokeCertificate(pid, int(reasonCode), sn); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func printJSON(data interface{}) (err error) {
	if err = encoder.Encode(data); err != nil {
		return err
	}
	return nil
}

const pwcharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$%&*-<>~"

func randomPassword(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = pwcharset[seededRand.Intn(len(pwcharset))]
	}
	return string(buf)
}
