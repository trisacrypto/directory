package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	db   store.Store
	conf config.Config
)

const dateFmt = "2006-01-02"

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "gdsutil"
	app.Version = pkg.Version()
	app.Usage = "utilities for managing a local gds instance"
	app.Commands = []*cli.Command{
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
			Name:     "vasp:list",
			Usage:    "list the VASPs in the current database by name, common name, and id",
			Category: "vasps",
			Action:   vaspList,
			Before:   connectDB,
			After:    closeDB,
			Flags:    []cli.Flag{},
		},
		{
			Name:      "vasp:detail",
			Usage:     "get the detail for a vasp record",
			ArgsUsage: "uuid [uuid ...]",
			Category:  "vasps",
			Action:    vaspDetail,
			Before:    connectDB,
			After:     closeDB,
			Flags:     []cli.Flag{},
		},
		{
			Name:     "vasp:status",
			Usage:    "inspect a VASP status and certificate requests",
			Category: "vasps",
			Action:   vaspStatus,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to get the status of",
					Required: true,
				},
			},
		},
		{
			Name:     "vasp:rereview",
			Usage:    "change a VASP verification status after it has been reviewed",
			Category: "vasps",
			Action:   rereview,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to rereview",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
				&cli.StringFlag{
					Name:    "verification-state",
					Aliases: []string{"state", "s"},
					Usage:   "specify the verification status for the VASP",
				},
			},
		},
		{
			Name:     "vasp:destroy",
			Usage:    "destroy a VASP record if it is in the rejected state",
			Category: "vasps",
			Action:   destroy,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to destroy the record for",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
			},
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
		{
			Name:     "contact:list",
			Usage:    "list the contacts in the current database",
			Category: "contact",
			Action:   contactList,
			Before:   connectDB,
			After:    closeDB,
			Flags:    []cli.Flag{},
		},
		{
			Name:     "contact:export",
			Usage:    "export the VASP contacts in the current database",
			Category: "contact",
			Action:   contactExport,
			Before:   connectDB,
			After:    closeDB,
			Flags:    []cli.Flag{},
		},
		{
			Name:      "contact:detail",
			Usage:     "get the detail for a contact record",
			ArgsUsage: "email [email ...]",
			Category:  "contact",
			Action:    contactDetail,
			Before:    connectDB,
			After:     closeDB,
			Flags:     []cli.Flag{},
		},
		{
			Name:     "certs:notify",
			Usage:    "send the reminder notification email that certs will be reissued soon",
			Category: "certs",
			Action:   notify,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to send reissuance reminder notifications to",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
				&cli.StringFlag{
					Name:    "reissuance-date",
					Aliases: []string{"d", "date"},
					Usage:   "the date that reissuance will occur in YYYY-MM-DD format",
					Value:   weekFromNow().Format(dateFmt),
				},
			},
		},
		{
			Name:     "certs:reissue",
			Usage:    "reissue identity certificates for the specified VASP",
			Category: "certs",
			Action:   reissueCerts,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to send reissuance reminder notifications to",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Usage:   "update the TRISA endpoint and common name for the new certs",
					Value:   "",
				},
				&cli.StringSliceFlag{
					Name:    "dns-names",
					Aliases: []string{"sans", "d"},
					Usage:   "specify additional DNS names to add to the request",
				},
				&cli.StringFlag{
					Name:    "webhook",
					Aliases: []string{"w"},
					Usage:   "specify a webhook to use to deliver certs",
				},
				&cli.BoolFlag{
					Name:    "no-email",
					Aliases: []string{"E"},
					Usage:   "do not deliver certificates by email",
				},
			},
		},
		{
			Name:     "certs:password",
			Usage:    "view or resend the password for the latest certificate request if available",
			Category: "certs",
			Action:   resendPassword,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to send reissuance reminder notifications to",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
				&cli.BoolFlag{
					Name:    "show",
					Aliases: []string{"s", "show-password"},
					Usage:   "show the password on the command line and exit without emailing the user",
					Value:   false,
				},
			},
		},
		{
			Name:     "certs:proto",
			Usage:    "create an identity certificate protocol buffer from a certificate",
			Category: "certs",
			Action:   makeCertificateProto,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "in",
					Aliases:  []string{"i"},
					Usage:    "path to identity certificates to convert on disk",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write json serialized identity certificate protocol buffers",
					Value:   "identity_certificate.pb.json",
				},
				&cli.StringFlag{
					Name:    "pkcs12password",
					Aliases: []string{"p"},
					Usage:   "pkcs12 password to decrypt certificates if required",
				},
			},
		},
		{
			Name:     "certs:revoke",
			Usage:    "mark a VASP as rejected and delete certificate information",
			Category: "certs",
			Action:   revokeCerts,
			Before:   connectDB,
			After:    closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to revoke the certificates of",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
			},
		},
		{
			Name:      "certs:dnsnames",
			Usage:     "add subject alternative names to the certificate request",
			ArgsUsage: "dnsName [dnsName ...]",
			Category:  "certs",
			Action:    addDNSNames,
			Before:    connectDB,
			After:     closeDB,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "vasp",
					Aliases:  []string{"vasp-id", "v"},
					Usage:    "the VASP ID to update certificate request records for",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "yes",
					Aliases: []string{"y"},
					Usage:   "skip the confirmation prompt and immediately send notifications",
					Value:   false,
				},
			},
		},

		{
			Name:     "certs:acme",
			Usage:    "verify a domain name via acme-dns challenge",
			Category: "certs",
			Action:   verifyDomain,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "token",
					Aliases: []string{"t"},
					Usage:   "if true, generate a challenge token for DNS verification",
				},
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "the domain to query the txt record for",
				},
				&cli.StringFlag{
					Name:    "challenge",
					Aliases: []string{"c"},
					Usage:   "the challenge token to verify",
				},
				&cli.BoolFlag{
					Name:    "debug",
					Aliases: []string{"D"},
					Usage:   "print the TXT records retrieved from DNS query",
				},
			},
		},
	}
	app.Run(os.Args)
}

//===========================================================================
// Before/After CLI Commands
//===========================================================================

func loadConf(*cli.Context) (err error) {
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

//===========================================================================
// Helper Functions
//===========================================================================

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func weekFromNow() time.Time {
	return time.Now().AddDate(0, 0, 7)
}

func printJSON(v interface{}) (err error) {
	if m, ok := v.(protoreflect.ProtoMessage); ok {
		return printJSONPB(m)
	}

	if msgs, ok := v.([]protoreflect.ProtoMessage); ok {
		return printJSONPBList(msgs)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(v); err != nil {
		return cli.Exit(fmt.Errorf("could not marshal json: %w", err), 1)
	}
	return nil
}

func printJSONPB(m protoreflect.ProtoMessage) (err error) {
	jsonpb := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = jsonpb.Marshal(m); err != nil {
		return cli.Exit(fmt.Errorf("could not marshal protocol buffer: %w", err), 1)
	}
	fmt.Println(string(data))
	return nil
}

func printJSONPBList(msgs []protoreflect.ProtoMessage) (err error) {
	objs := make([]map[string]interface{}, len(msgs))
	for i, msg := range msgs {
		if objs[i], err = wire.Rewire(msg); err != nil {
			return cli.Exit(fmt.Errorf("could not rewire message %d: %w", i, err), 1)
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(objs); err != nil {
		return cli.Exit(fmt.Errorf("could not marshal json: %w", err), 1)
	}
	return nil
}
