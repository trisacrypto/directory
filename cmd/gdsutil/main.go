package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"github.com/trisacrypto/directory/pkg"
	profiles "github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	storerr "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
)

var (
	profile *profiles.Profile
	db      store.Store
	conf    config.Config
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "gdsutil"
	app.Version = pkg.Version()
	app.Usage = "utilities for operating the GDS service and database"
	app.Before = loadProfile
	app.Commands = []*cli.Command{
		{
			Name:      "profile",
			Aliases:   []string{"config", "profiles"},
			Usage:     "view and manage profiles to configure gdsutil with",
			UsageText: "gdsutil profile [name]\n   gdsutil profile --activate [name]\n   gdsutil profile --list\n   gdsutil profile --path\n   gdsutil profile --install\n   gdsutil profile --edit",
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
		{
			Name:      "decrypt",
			Usage:     "decrypt base64 encoded ciphertext with an HMAC signature",
			ArgsUsage: "ciphertext hmac",
			Category:  "cipher",
			Action:    cipherDecrypt,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "key",
					Aliases: []string{"k"},
					Usage:   "secret key to decrypt the cipher text",
					EnvVars: []string{"GDS_SECRET_KEY"},
				},
			},
		},
		{
			Name:     "admin:tokenkey",
			Usage:    "generate an RSA token key pair and ksuid for JWT token signing",
			Category: "admin",
			Action:   generateTokenKey,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write keys out to (optional, will be saved as ksuid.pem by default)",
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "number of bits for the generated keys",
					Value:   4096,
				},
			},
		},
		{
			Name:     "db:usage",
			Usage:    "count the number of objects in the database by namespace",
			Category: "db",
			Action:   dbUsage,
			Before:   connectDB,
			After:    closeDB,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "contact:migrate",
			Usage:    "migrate all contacts on vasps into the model contacts namespace",
			Category: "contact",
			Action:   migrateContacts,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dryrun",
					Aliases: []string{"d"},
					Usage:   "print migration results without modifying the database, used for testing",
					Value:   true,
				},
				&cli.BoolFlag{
					Name:    "compare",
					Aliases: []string{"c"},
					Usage:   "if the contact exists, compare to the vasp contact record",
				},
			},
		},
		{
			Name:     "contact:fixverifytoken",
			Usage:    "fixes any unverified contacts that do not have verification tokens",
			Category: "contact",
			Action:   fixVerifyToken,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "dryrun",
					Aliases: []string{"d"},
					Usage:   "print migration results without modifying the database, used for testing",
				},
			},
		},
	}
	app.Run(os.Args)
}

//===========================================================================
// Profile Actions
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
// Cipher Actions
//===========================================================================

const nonceSize = 12

func cipherDecrypt(c *cli.Context) (err error) {
	if c.NArg() != 2 {
		return cli.Exit("must specify ciphertext and hmac arguments", 1)
	}

	var secret string
	if secret = c.String("key"); secret == "" {
		return cli.Exit("cipher key required", 1)
	}

	var ciphertext, signature []byte
	if ciphertext, err = base64.RawStdEncoding.DecodeString(c.Args().Get(0)); err != nil {
		return cli.Exit(fmt.Errorf("could not decode ciphertext: %s", err), 1)
	}
	if signature, err = base64.RawStdEncoding.DecodeString(c.Args().Get(1)); err != nil {
		return cli.Exit(fmt.Errorf("could not decode signature: %s", err), 1)
	}

	if len(ciphertext) == 0 {
		return cli.Exit("empty cipher text", 1)
	}

	// Create a 32 byte signature of the key
	hash := sha256.New()
	hash.Write([]byte(secret))
	key := hash.Sum(nil)

	// Separate the data from the nonce
	data := ciphertext[:len(ciphertext)-nonceSize]
	nonce := ciphertext[len(ciphertext)-nonceSize:]

	// Validate HMAC signature
	if err = validateHMAC(key, data, signature); err != nil {
		return cli.Exit(err, 1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return cli.Exit(err, 1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return cli.Exit(err, 1)
	}

	plainbytes, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(plainbytes))
	return nil
}

//===========================================================================
// Cipher Helper Functions
//===========================================================================

func createHMAC(key, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot sign empty data")
	}
	hm := hmac.New(sha256.New, key)
	hm.Write(data)
	return hm.Sum(nil), nil
}

func validateHMAC(key, data, sig []byte) error {
	hmac, err := createHMAC(key, data)
	if err != nil {
		return err
	}

	if !bytes.Equal(sig, hmac) {
		return errors.New("HMAC mismatch")
	}
	return nil
}

//===========================================================================
// Admin Functions
//===========================================================================

func generateTokenKey(c *cli.Context) (err error) {
	// Create ksuid and determine outpath
	var keyid ksuid.KSUID
	if keyid, err = ksuid.NewRandom(); err != nil {
		return cli.Exit(err, 1)
	}

	var out string
	if out = c.String("out"); out == "" {
		out = fmt.Sprintf("%s.pem", keyid)
	}

	// Generate RSA keys using crypto random
	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, c.Int("size")); err != nil {
		return cli.Exit(err, 1)
	}

	// Open file to PEM encode keys to
	var f *os.File
	if f, err = os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return cli.Exit(err, 1)
	}

	if err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("RSA key id: %s -- saved with PEM encoding to %s\n", keyid, out)
	return nil
}

//===========================================================================
// Database Functions
//===========================================================================

func dbUsage(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	counters := []struct {
		namespace string
		count     func(context.Context) (uint64, error)
	}{
		{wire.NamespaceVASPs, db.CountVASPs},
		{wire.NamespaceCertReqs, db.CountCertReqs},
		{wire.NamespaceCerts, db.CountCerts},
		{wire.NamespaceAnnouncements, db.CountAnnouncementMonths},
		{wire.NamespaceActivities, db.CountActivityMonth},
		{wire.NamespaceOrganizations, db.CountOrganizations},
		{wire.NamespaceContacts, db.CountContacts},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.AlignRight)
	for _, counter := range counters {
		var count uint64
		if count, err = counter.count(ctx); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Fprintf(w, "%s\t%d\n", counter.namespace, count)
	}
	w.Flush()
	return nil
}

//===========================================================================
// Contact Functions
//===========================================================================

func migrateContacts(c *cli.Context) (err error) {
	dryrun := c.Bool("dryrun")
	compare := c.Bool("compare")

	// Iterate through all vasps in the database
	vasps := db.ListVASPs(context.Background())
	defer vasps.Release()
	for vasps.Next() {
		var vasp *pb.VASP
		if vasp, err = vasps.VASP(); err != nil {
			return cli.Exit(err, 1)
		}

		// Iterate through all contacts on the vasp
		contacts := models.NewContactIterator(vasp.Contacts, models.SkipNoEmail())
		for contacts.Next() {
			vaspContact, _ := contacts.Value()

			var contact *models.Contact
			if contact, err = db.RetrieveContact(context.Background(), vaspContact.Email); err != nil {
				if errors.Is(err, storerr.ErrEntityNotFound) {
					if dryrun {
						fmt.Printf("contact %s missing and needs to be migrated\n\n", vaspContact.Email)
					} else {
						// Create the contact if it doesn't exist
						extra := &models.GDSContactExtraData{}
						if err := vaspContact.Extra.UnmarshalTo(extra); err != nil {
							return cli.Exit(fmt.Errorf("could not unmarshal extra for %s: %s", vaspContact.Email, err), 1)
						}

						contact = &models.Contact{
							Email:    vaspContact.Email,
							Name:     vaspContact.Name,
							Vasps:    []string{vasp.CommonName},
							Verified: extra.Verified,
							Token:    extra.Token,
							EmailLog: extra.EmailLog,
							Created:  time.Now().Format(time.RFC3339),
							Modified: time.Now().Format(time.RFC3339),
						}

						if _, err := db.CreateContact(context.Background(), contact); err != nil {
							return cli.Exit(err, 1)
						}
						fmt.Printf("contact %s created!\n", vaspContact.Email)
					}
					continue
				}
				return cli.Exit(err, 1)
			}

			if !compare {
				continue
			}

			// Check contact equality
			updates := make(map[string]map[string]interface{})
			if vaspContact.Email != contact.Email {
				updates["email"] = map[string]interface{}{
					"contact": contact.Email,
					"vasp":    vaspContact.Email,
				}
			}

			if vaspContact.Name != contact.Name {
				updates["name"] = map[string]interface{}{
					"contact": contact.Name,
					"vasp":    vaspContact.Name,
				}
			}

			found := false
			for _, included := range contact.Vasps {
				if vasp.CommonName == included {
					found = true
					break
				}
			}

			if !found {
				updates["vasps"] = map[string]interface{}{
					"contact": contact.Vasps,
					"vasp":    append(contact.Vasps, vasp.CommonName),
				}
			}

			extra := &models.GDSContactExtraData{}
			if err := vaspContact.Extra.UnmarshalTo(extra); err != nil {
				return cli.Exit(fmt.Errorf("could not unmarshal extra for %s: %s", vaspContact.Email, err), 1)
			}

			if extra.Verified != contact.Verified {
				updates["verified"] = map[string]interface{}{
					"contact": contact.Verified,
					"vasp":    extra.Verified,
				}
			}

			if extra.Token != contact.Token {
				updates["token"] = map[string]interface{}{
					"contact": contact.Token,
					"vasp":    extra.Token,
				}
			}

			if !slices.Equal(extra.EmailLog, contact.EmailLog) {
				updates["email_log"] = map[string]interface{}{
					"contact": len(contact.EmailLog),
					"vasp":    len(extra.EmailLog),
				}
			}

			if len(updates) > 0 {
				if dryrun {
					fmt.Printf("contact %s does not match vasp contact and needs to be updated\n", contact.Email)
					w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', tabwriter.AlignRight)
					fmt.Fprintln(w, "field\tcontact\tvasp")
					for field, vals := range updates {
						fmt.Fprintf(w, "%s\t%v\t%v\n", field, vals["contact"], vals["vasp"])
					}
					w.Flush()
					fmt.Println()
				} else {
					fmt.Printf("contact %s updated!\n", contact.Email)
				}
			}
		}
	}

	if err = vasps.Error(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func fixVerifyToken(c *cli.Context) (err error) {
	dryrun := c.Bool("dryrun")

	vasps := db.ListVASPs(context.Background())
	defer vasps.Release()
	for vasps.Next() {
		var vasp *pb.VASP
		if vasp, err = vasps.VASP(); err != nil {
			return cli.Exit(err, 1)
		}

		contacts := models.NewContactIterator(vasp.Contacts, models.SkipNoEmail(), models.SkipDuplicates())
		for contacts.Next() {
			vaspContact, _ := contacts.Value()

			var contact *models.Contact
			if contact, err = db.RetrieveContact(context.Background(), vaspContact.Email); err != nil {
				return cli.Exit(err, 1)
			}

			if !contact.Verified && contact.Token == "" {
				fmt.Printf("contact %s is not verified and has no verification token\n", contact.Email)

				if !dryrun {
					contact.Token = secrets.CreateToken(models.VerificationTokenLength)
					if err = db.UpdateContact(context.Background(), contact); err != nil {
						return cli.Exit(err, 1)
					}
				}
			}
		}
	}

	if err = vasps.Error(); err != nil {
		return cli.Exit(err, 1)
	}

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

func connectDB(c *cli.Context) (err error) {
	// Suppress the zerolog output from the store.
	logger.Discard()

	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	conf.Database.ReindexOnBoot = false
	conf.ConsoleLog = false

	if db, err = store.Open(conf.Database); err != nil {
		if serr, ok := status.FromError(err); ok {
			return cli.Exit(fmt.Errorf("could not open store: %s", serr.Message()), 1)
		}
		return cli.Exit(err, 1)
	}

	return nil
}

func closeDB(c *cli.Context) (err error) {
	if db != nil {
		if err = db.Close(); err != nil {
			return cli.Exit(err, 2)
		}
	}
	return nil
}
