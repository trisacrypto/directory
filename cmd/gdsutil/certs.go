package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/trisacrypto/courier/pkg/api/v1"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/whisper"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trust"

	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

//===========================================================================
// Certificate Management
//===========================================================================

const (
	whisperPasswordTemplate = "Below is the PKCS12 password which you must use to decrypt your new certificates:\n\n%s\n"
)

func notify(c *cli.Context) (err error) {
	var (
		nsent       int
		vasp        *pb.VASP
		emailer     *emails.EmailManager
		reissueDate time.Time
	)

	vaspID := c.String("vasp")
	fmt.Printf("looking up vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Step 1: Fetch the VASP record
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
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

	fmt.Printf("successfully sent %d email reminders about certificate reissuance on %s\n", nsent, reissueDate.Format(emails.DateFormat))
	return nil
}

func reissueCerts(c *cli.Context) (err error) {
	var (
		vasp           *pb.VASP
		certreq        *models.CertificateRequest
		pkcs12password string
		emailer        *emails.EmailManager
		whisperLink    string
		nsent          int
	)

	vaspID := c.String("vasp")
	fmt.Printf("looking up vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Step 1: Fetch the VASP record
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	// Check with the user if we should continue with the certificate reissuance
	fmt.Printf("reissuing certs for %s\n", vasp.CommonName)
	if !c.Bool("yes") {
		if !askForConfirmation("continue with certificate reissuance?") {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	if endpoint := c.String("endpoint"); endpoint != "" {
		var cname string
		if cname, _, err = net.SplitHostPort(endpoint); err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Printf("updating common name to %s and endpoint to %s\n", cname, endpoint)
		vasp.CommonName = cname
		vasp.TrisaEndpoint = endpoint
	}

	// Step 2: Create a CertificateRequest
	if certreq, err = models.NewCertificateRequest(vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not create certificate request: %s", err), 1)
	}

	// Step 2b: add any additional dns names from the command line
	if dnsNames := c.StringSlice("dns-names"); len(dnsNames) > 0 {
		certreq.DnsNames = append(certreq.DnsNames, dnsNames...)
	}

	// Override the certificate delivery webhook if specified
	if webhook := c.String("webhook"); webhook != "" {
		certreq.Webhook = webhook
	}

	// Override the email delivery preference
	certreq.NoEmailDelivery = c.Bool("no-email")

	// Step 2c: mark the certificate request as ready to submit for CertMan
	if err = models.UpdateCertificateRequestStatus(
		certreq,
		models.CertificateRequestState_READY_TO_SUBMIT,
		"manually reissuing certificates",
		"support@rotational.io",
	); err != nil {
		return cli.Exit(fmt.Errorf("could not mark certificate request ready to submit: %s", err), 1)
	}

	// Step 3: Create a PKCS12 password
	var sm *secrets.SecretManager
	if sm, err = secrets.New(conf.Secrets); err != nil {
		return cli.Exit(fmt.Errorf("could not connect to secret manager: %s", err), 1)
	}

	secretType := "password"
	pkcs12password = secrets.CreateToken(16)

	if err = sm.With(certreq.Id).CreateSecret(ctx, secretType); err != nil {
		return cli.Exit(fmt.Errorf("could not create password secret: %s", err), 1)
	}

	if err = sm.With(certreq.Id).AddSecretVersion(ctx, secretType, []byte(pkcs12password)); err != nil {
		return cli.Exit(fmt.Errorf("could not create password version: %s", err), 1)
	}

	if !certreq.NoEmailDelivery {
		// Create a Whisper link for the provided PKCS12 password.
		if whisperLink, err = whisper.CreateSecretLink(fmt.Sprintf(whisperPasswordTemplate, pkcs12password), "", 3, weekFromNow()); err != nil {
			return cli.Exit(err, 1)
		}

		// Create the email manager.
		if emailer, err = emails.New(conf.Email); err != nil {
			return cli.Exit(err, 1)
		}

		// Send the notification email that certificate reissuance is forthcoming and provide whisper link to the PKCS12 password.
		if nsent, err = emailer.SendReissuanceStarted(vasp, whisperLink); err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Printf("successfully sent %d Whisper password notifications for PKCS12 password %q\n", nsent, pkcs12password)
	}

	if certreq.Webhook != "" {
		// Create a courier client to deliver the pkcs12 password to the TRISA member
		var client api.CourierClient
		if client, err = api.New(certreq.Webhook); err != nil {
			return cli.Exit(fmt.Errorf("could not create courier client: %s", err), 1)
		}

		// Deliver the pkcs12 password via the webhook
		req := &api.StorePasswordRequest{
			ID:       certreq.Id,
			Password: pkcs12password,
		}
		if err = client.StoreCertificatePassword(ctx, req); err != nil {
			return cli.Exit(fmt.Errorf("could not deliver pkcs12 password with webhook: %s", err), 1)
		}

		fmt.Printf("successfully sent PKCS12 password to webhook %q\n", certreq.Webhook)
	}

	// Save certificate request to database
	if err = db.UpdateCertReq(ctx, certreq); err != nil {
		return cli.Exit(fmt.Errorf("could not save certreq: %s", err), 1)
	}

	// Step 4: Save certificate request on VASP
	if err = models.AppendCertReqID(vasp, certreq.Id); err != nil {
		return cli.Exit(fmt.Errorf("could not append certreq to VASP: %s", err), 1)
	}

	if err = db.UpdateVASP(ctx, vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not save vasp: %s", err), 1)
	}

	return nil
}

func cancelCertificatRequest(c *cli.Context) (err error) {
	var (
		vasp    *pb.VASP
		certreq *models.CertificateRequest
	)

	reqID := c.String("request")
	if reqID == "" {
		return cli.Exit("specify a certificate request ID to cancel", 1)
	}

	ctx := context.Background()
	if certreq, err = db.RetrieveCertReq(ctx, reqID); err != nil {
		return cli.Exit(fmt.Errorf("could not find certificate rqeuest: %w", err), 1)
	}

	// Check with the user if we should continue with canceling the request
	fmt.Printf("found certificate request for %s status %s\n", certreq.CommonName, certreq.Status)
	if !c.Bool("yes") {
		if !askForConfirmation("cancel this request?") {
			return cli.Exit(fmt.Errorf("operation halted by user"), 1)
		}
	}

	// Delete the certificate request
	if err = db.DeleteCertReq(ctx, reqID); err != nil {
		return cli.Exit(fmt.Errorf("could not delete certificate request: %w", err), 1)
	}

	// Remove the certificate request record from the vasp
	if vasp, err = db.RetrieveVASP(ctx, certreq.Vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not retrieve vasp: %w", err), 1)
	}

	if err = models.DeleteCertReqID(vasp, reqID); err != nil {
		return cli.Exit(fmt.Errorf("could not update vasp: %w", err), 1)
	}

	// Change the status of the VASP
	if !c.Bool("no-status-change") {
		if vasp.VerificationStatus == pb.VerificationIssuing || vasp.VerificationStatus == pb.VerificationReviewed {
			fmt.Printf("updating status of %s to %s\n", vasp.CommonName, pb.VerificationPending)
			if !c.Bool("yes") {
				if !askForConfirmation("continue?") {
					return cli.Exit(fmt.Errorf("operation halted by user"), 1)
				}
			}

			vasp.VerificationStatus = pb.VerificationPending
		}
	}

	if err = db.UpdateVASP(ctx, vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not update vasp: %w", err), 1)
	}

	return nil
}

func resendPassword(c *cli.Context) (err error) {
	var (
		vasp           *pb.VASP
		vaspName       string
		certreqID      string
		pkcs12password []byte
		sm             *secrets.SecretManager
		emailer        *emails.EmailManager
		whisperLink    string
		nsent          int
	)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Fetch and identify the VASP specified by the user
	vaspID := c.String("vasp")
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record %s: %s", vaspID, err), 1)
	}

	if vaspName, err = vasp.Name(); err != nil {
		vaspName = vasp.CommonName
	}

	// Get the latest certificate request for the VASP
	if certreqID, err = models.GetLatestCertReqID(vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not get latest certificate request ID for vasp %s: %s", vaspName, err), 1)
	}

	// Connect to the secrets store and fetch the PKCS12 password if it exists
	if sm, err = secrets.New(conf.Secrets); err != nil {
		return cli.Exit(fmt.Errorf("could not connect to secret manager: %s", err), 1)
	}

	if pkcs12password, err = sm.With(certreqID).GetLatestVersion(ctx, "password"); err != nil {
		return cli.Exit(fmt.Errorf("could not retrieve pkcs12 password for vasp %s certificate request %s: %s", vaspName, certreqID, err), 1)
	}

	// If print password and exit, do that without user confirmation
	if c.Bool("show") {
		fmt.Printf("retrieved password for %s (certificate request %s)\n", vaspName, certreqID)
		if !c.Bool("yes") {
			if !askForConfirmation("show PKCS12 password on the command line?") {
				return cli.Exit(fmt.Errorf("canceled by user"), 1)
			}
		}

		// Print password and exit
		fmt.Println(string(pkcs12password))
		return nil
	}

	// Check with the user if we should continue with resending the password
	fmt.Printf("resending password for %s (certificate request %s)\n", vaspName, certreqID)
	if !c.Bool("yes") {
		if !askForConfirmation("continue and resend PKCS12 password?") {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	// Create the Whisper link for the provided PKCS12 password.
	if whisperLink, err = whisper.CreateSecretLink(fmt.Sprintf(whisperPasswordTemplate, string(pkcs12password)), "", 3, weekFromNow()); err != nil {
		return cli.Exit(err, 1)
	}

	// Create the email manager.
	if emailer, err = emails.New(conf.Email); err != nil {
		return cli.Exit(err, 1)
	}

	// Send the notification email that certificate reissuance is forthcoming and provide whisper link to the PKCS12 password.
	if nsent, err = emailer.SendReissuanceStarted(vasp, whisperLink); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("successfully sent %d Whisper password notifications for PKCS12 password %q\n", nsent, pkcs12password)
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

	if err = os.WriteFile(c.String("out"), data, 0600); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func revokeCerts(c *cli.Context) (err error) {
	vaspID := c.String("vasp")
	fmt.Printf("lookup vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	var vasp *pb.VASP
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	// Check with the user if we should continue with the certificate revocation
	fmt.Printf("revoking certs for %s\n", vasp.CommonName)
	if !c.Bool("yes") {
		if !askForConfirmation("continue with certificate revocation?") {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	// Mark any outstanding certificate requests as rejected.
	var certreqs []string
	if certreqs, err = models.GetCertReqIDs(vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not get certificate request ids: %s", err), 1)
	}

	for _, crid := range certreqs {
		var certreq *models.CertificateRequest
		if certreq, err = db.RetrieveCertReq(ctx, crid); err != nil {
			fmt.Printf("error retrieving certreq %s: %s\n", crid, err)
			continue
		}

		// Only update certreqs that are not completed.
		if certreq.Status < models.CertificateRequestState_COMPLETED {
			if err = models.UpdateCertificateRequestStatus(
				certreq,
				models.CertificateRequestState_CR_REJECTED,
				"certificate request canceled by admin",
				"support@rotational.io",
			); err != nil {
				fmt.Printf("could not mark certificate request %s as rejected: %s\n", crid, err)
				continue
			}

			if err = db.UpdateCertReq(ctx, certreq); err != nil {
				fmt.Printf("could not save certreq %s: %s\n", crid, err)
			}

			fmt.Printf("marked certificate request %s as rejected", crid)
		}
	}

	// Print the serial number for revocation in Sectigo
	if vasp.IdentityCertificate != nil {
		fmt.Printf("please revoke certificates with serial number %X\n", vasp.IdentityCertificate.SerialNumber)
	}

	if len(vasp.SigningCertificates) > 0 {
		for _, cert := range vasp.SigningCertificates {
			fmt.Printf("please revoke certificates with serial number %X\n", cert.SerialNumber)
		}
	}

	// TODO: what do we have to do with the certificates models?
	vasp.IdentityCertificate = nil
	vasp.SigningCertificates = nil

	// Set the VASP state to rejected
	if err = models.UpdateVerificationStatus(
		vasp,
		pb.VerificationState_REJECTED,
		"certificates revoked due to cessation of operations",
		"support@rotational.io",
	); err != nil {
		return cli.Exit(fmt.Errorf("could not update VASP status: %s", err), 1)
	}

	if err = db.UpdateVASP(ctx, vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not save VASP: %s", err), 1)
	}

	fmt.Println("VASP registration revoked")
	return nil
}

func verifyDomain(c *cli.Context) (err error) {
	if token := c.Bool("token"); token {
		// Generate a challenge token, print and return
		// TODO: also generate whisper link to send to user
		nonce := make([]byte, 32)
		if _, err = rand.Read(nonce); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Printf("Challenge Token: TRISA-DOMAIN-VERIFICATION=%s\n", base64.RawURLEncoding.EncodeToString(nonce))
		return nil
	}

	var (
		domain    string
		challenge string
		debug     bool
	)

	debug = c.Bool("debug")

	if domain = c.String("domain"); domain == "" {
		return cli.Exit("domain required for challenge", 1)
	}

	if challenge = c.String("challenge"); challenge == "" {
		return cli.Exit("challenge required for domain verification", 1)
	}

	var answers []string
	if answers, err = net.LookupTXT(domain); err != nil {
		return cli.Exit(err, 1)
	}

	challenge = strings.TrimSpace(challenge)
	for _, answer := range answers {
		if debug {
			fmt.Println(answer)
		}

		if strings.TrimSpace(answer) == challenge {
			fmt.Println("domain verified!")
			return nil
		}
	}

	return cli.Exit(fmt.Errorf("%d TXT records returned did not match challenge", len(answers)), 1)
}

func addDNSNames(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.Exit("specify at least one dns name to add", 1)
	}

	var dnsNames []string
	for i := 0; i < c.NArg(); i++ {
		if name := c.Args().Get(i); name != "" {
			dnsNames = append(dnsNames, name)
		}
	}

	vaspID := c.String("vasp")
	fmt.Printf("lookup vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	var vasp *pb.VASP
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	certreqs, err := models.GetCertReqIDs(vasp)
	if err != nil {
		return cli.Exit(err, 1)
	}

	for _, certreq := range certreqs {
		ca, err := db.RetrieveCertReq(ctx, certreq)
		if err != nil {
			return cli.Exit(err, 1)
		}

		if ca.Status <= models.CertificateRequestState_READY_TO_SUBMIT {
			// Check with the user if we should continue with the certificate revocation
			fmt.Printf("updating certificate requests for %s\n", vasp.CommonName)
			if !c.Bool("yes") {
				if !askForConfirmation(fmt.Sprintf("add %d dns names to %q?", len(dnsNames), ca.BatchName)) {
					return cli.Exit(fmt.Errorf("canceled by user"), 1)
				}
			}

			ca.DnsNames = append(ca.DnsNames, dnsNames...)

			if err = db.UpdateCertReq(ctx, ca); err != nil {
				return cli.Exit(err, 1)
			}
		}
	}
	return nil
}
