/*
Reissuer is a quick command to help us easily reissue expiring certificates and send
notification emails manually. This is a stopgap solution to automated reissuance and
this command will eventually be migrated to the gdsutil.
*/
package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trust"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	db   store.Store
	conf config.Config
)

const (
	dateFmt = "2006-01-02"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "reissuer"
	app.Version = pkg.Version()
	app.Usage = "a quick helper tool to manually reissue expiring certs"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:   "notify",
			Usage:  "send the reminder notification email that certs will be reissued soon",
			Action: notify,
			Before: connectDB,
			After:  closeDB,
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
			Name:   "certs",
			Usage:  "reissue identity certificates for the specified VASP",
			Action: reissueCerts,
			Before: connectDB,
			After:  closeDB,
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
			},
		},
		{
			Name:   "proto",
			Usage:  "create an identity certificate protocol buffer from a certificate",
			Action: makeCertificateProto,
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
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func connectDB(c *cli.Context) (err error) {
	// Surpress the zerolog output from the store
	logger.Discard()

	// Load the configuration from the environment
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	conf.Database.ReindexOnBoot = false
	conf.ConsoleLog = false

	// Connect to the trtl server and create a store to access data directly like GDS
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
		return cli.Exit(err, 2)
	}
	return nil
}

func notify(c *cli.Context) (err error) {
	var (
		nsent       int
		vasp        *pb.VASP
		emailer     *emails.EmailManager
		reissueDate time.Time
	)

	vaspID := c.String("vasp")
	fmt.Printf("looking up vasp with id %s\n", vaspID)

	// Step 1: Fetch the VASP record
	if vasp, err = db.RetrieveVASP(vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	// Check with the user if we should continue to send reissuance reminder emails
	fmt.Printf("reissuing certs for %s\n", vasp.CommonName)
	if !c.Bool("yes") {
		if !askForConfirmation("continue sending reissuance reminder emails?") {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	// Step 2: Parse reissuance date or get date 1 week from today
	if reissueDate, err = time.Parse(dateFmt, c.String("reissuance-date")); err != nil {
		return cli.Exit(err, 1)
	}

	// Step 3: Connect to the Email Manager with SendGrid API client and send emails
	if emailer, err = emails.New(conf.Email); err != nil {
		return cli.Exit(err, 1)
	}

	// Send reissuance reminder emails
	if nsent, err = emailer.SendReissuanceReminder(vasp, reissueDate); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("successfully sent %d email remdiners about certificate reissuance on %s\n", nsent, reissueDate.Format(emails.DateFormat))
	return nil
}

func reissueCerts(c *cli.Context) (err error) {
	var (
		vasp           *pb.VASP
		certreq        *models.CertificateRequest
		pkcs12password string
	)

	vaspID := c.String("vasp")
	fmt.Printf("looking up vasp with id %s\n", vaspID)

	// Step 1: Fetch the VASP record
	if vasp, err = db.RetrieveVASP(vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	// Check with the user if we should continue with the certificate reissuance
	fmt.Printf("reissuing certs for %s\n", vasp.CommonName)
	if !c.Bool("yes") {
		if !askForConfirmation("continue with certificate reissuance?") {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	// Step 2: Create a CertificateRequest
	if certreq, err = models.NewCertificateRequest(vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not create certificate request: %s", err), 1)
	}

	// Step 2b: mark the certificate request as ready to submit for CertMan
	if err = models.UpdateCertificateRequestStatus(
		certreq,
		models.CertificateRequestState_READY_TO_SUBMIT,
		"manually reissuing certificates",
		"support@rotational.io",
	); err != nil {
		return cli.Exit(fmt.Errorf("could not mark certificate request ready to submit: %s", err), 1)
	}

	// Step 3: Create a PKCS12 password and print it out
	var sm *secrets.SecretManager
	if sm, err = secrets.New(conf.Secrets); err != nil {
		return cli.Exit(fmt.Errorf("could not connect to secret manager: %s", err), 1)
	}

	secretType := "password"
	pkcs12password = secrets.CreateToken(16)

	// TODO: instead of printing password create whisper link and send email with link.
	fmt.Printf("PKCS12 Password is: %q\n", pkcs12password)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = sm.With(certreq.Id).CreateSecret(ctx, secretType); err != nil {
		return cli.Exit(fmt.Errorf("could not create password secret: %s", err), 1)
	}
	if err = sm.With(certreq.Id).AddSecretVersion(ctx, secretType, []byte(pkcs12password)); err != nil {
		return cli.Exit(fmt.Errorf("could not create password version: %s", err), 1)
	}

	// Save certificate request to database
	if err = db.UpdateCertReq(certreq); err != nil {
		return cli.Exit(fmt.Errorf("could not save certreq: %s", err), 1)
	}

	// Step 4: Save certificate request on VASP
	if err = models.AppendCertReqID(vasp, certreq.Id); err != nil {
		return cli.Exit(fmt.Errorf("could not append certreq to VASP: %s", err), 1)
	}

	if err = db.UpdateVASP(vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not save vasp: %s", err), 1)
	}

	return nil
}

func makeCertificateProto(c *cli.Context) (err error) {
	var archive *trust.Serializer
	if pkcs12password := c.String("pkcs12password"); pkcs12password != "" {
		if archive, err = trust.NewSerializer(true, pkcs12password, trust.CompressionAuto); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		if archive, err = trust.NewSerializer(false); err != nil {
			return cli.Exit(err, 1)
		}
	}

	var provider *trust.Provider
	if provider, err = archive.ReadFile(c.String("in")); err != nil {
		return cli.Exit(err, 1)
	}

	var cert *x509.Certificate
	if cert, err = provider.GetLeafCertificate(); err != nil {
		return cli.Exit(err, 1)
	}

	pub := &pb.Certificate{
		Version:            int64(cert.Version),
		SerialNumber:       cert.SerialNumber.Bytes(),
		Signature:          cert.Signature,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		Subject: &pb.Name{
			CommonName:         cert.Subject.CommonName,
			SerialNumber:       cert.Subject.SerialNumber,
			Organization:       cert.Subject.Organization,
			OrganizationalUnit: cert.Subject.OrganizationalUnit,
			StreetAddress:      cert.Subject.StreetAddress,
			Locality:           cert.Subject.Locality,
			Province:           cert.Subject.Province,
			PostalCode:         cert.Subject.PostalCode,
			Country:            cert.Subject.Country,
		},
		Issuer: &pb.Name{
			CommonName:         cert.Issuer.CommonName,
			SerialNumber:       cert.Issuer.SerialNumber,
			Organization:       cert.Issuer.Organization,
			OrganizationalUnit: cert.Issuer.OrganizationalUnit,
			StreetAddress:      cert.Issuer.StreetAddress,
			Locality:           cert.Issuer.Locality,
			Province:           cert.Issuer.Province,
			PostalCode:         cert.Issuer.PostalCode,
			Country:            cert.Issuer.Country,
		},
		NotBefore: cert.NotBefore.Format(time.RFC3339),
		NotAfter:  cert.NotAfter.Format(time.RFC3339),
		Revoked:   false,
	}

	// Write the public certificate into the directory service data store
	buf := bytes.NewBuffer(nil)
	if err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		return cli.Exit(err, 1)
	}
	pub.Data = buf.Bytes()

	// Write the entire provider chain into the directory service data store
	if archive, err = trust.NewSerializer(false, "", trust.CompressionGZIP); err != nil {
		return cli.Exit(err, 1)
	}

	// Ensure only the public keys are written to the directory service
	if pub.Chain, err = archive.Compress(provider.Public()); err != nil {
		return cli.Exit(err, 1)
	}

	// Now serialize and base64 encode the certificate
	var data []byte
	if data, err = protojson.Marshal(pub); err != nil {
		return cli.Exit(err, 1)
	}

	if err = ioutil.WriteFile(c.String("out"), data, 0600); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

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
