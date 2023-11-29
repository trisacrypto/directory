package certman

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	courier "github.com/trisacrypto/courier/pkg/api/v1"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/whisper"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func New(conf config.CertManConfig, db store.Store, secret *secrets.SecretManager, email *emails.EmailManager) (_ Service, err error) {
	// If not enabled return the certman disabled stub.
	if !conf.Enabled {
		return &Disabled{}, nil
	}

	// If enabled, construct the certificate manager and return it
	cm := &CertificateManager{
		conf:   conf,
		db:     db,
		secret: secret,
		email:  email,
	}

	if cm.email == nil {
		return nil, errors.New("email manager is required for cert manager")
	}

	if cm.secret == nil {
		return nil, errors.New("secret manager is required for cert manager")
	}

	if conf.Sectigo.Testing() {
		if err = mock.Start(conf.Sectigo.Profile); err != nil {
			return nil, err
		}
	}

	if cm.certs, err = sectigo.New(conf.Sectigo); err != nil {
		return nil, err
	}

	if cm.certDir, err = cm.getCertStorage(); err != nil {
		return nil, err
	}

	return cm, nil
}

// CertificateManager is a struct with a go routine that periodically checks on the
// status of certificate requests and moves them through the request pipeline. This is
// separated from the parent GDS to allow for isolated testing.
type CertificateManager struct {
	conf    config.CertManConfig
	db      store.Store
	secret  *secrets.SecretManager
	certs   *sectigo.Sectigo
	email   *emails.EmailManager
	certDir string
	stop    chan struct{}
}

// Compile time interface implementation check.
var _ Service = &CertificateManager{}

// Run starts the CertManager as a go routine under the provided waitgroup. For
// graceful shutdown, the caller must invoke the Stop method to signal the CertManager
// routine to stop and block on the waitgroup if provided.
func (c *CertificateManager) Run(wg *sync.WaitGroup) error {
	if !c.conf.Enabled {
		return errors.New("certificate manager is not enabled")
	}

	if c.stop != nil {
		return errors.New("certificate manager is already running")
	}

	if wg != nil {
		wg.Add(1)
	}

	c.stop = make(chan struct{})
	go func() {
		c.CertManager()
		c.stop = nil
		if wg != nil {
			wg.Done()
		}
	}()
	return nil
}

// Stop signals the CertManager routine to shutdown.
// Note: This does not wait for the CertManager to stop and the caller should block on
// the waitgroup passed to the Run method in order to implement a graceful shutdown.
func (c *CertificateManager) Stop() {
	if c.stop != nil {
		close(c.stop)
	}
}

// CertManager is a go routine that periodically checks on the status of certificate
// requests and moves them through the request pipeline. Once CertManager detects a
// certificate request that is ready to submit, it submits the request via the Sectigo
// API. If processing, it checks the batch status, and when it detects that the bact
// is done processing it downloads the certs and emails them to the technical contacts.
// If the certificate processing fails for any reason, it sends an error message to
// the TRISA admins since this will prevent the integrator from joining the network.
//
// TODO: move completed certificate requests to archive so that the CertManger routine
// isn't continuously handling a growing number of requests over time.
func (c *CertificateManager) CertManager() {
	// Tickers are created in the go routine to prevent backpressure if the individual
	// handler routines take longer than the ticker intervals.
	// Note: These routines block in order to facilitate a graceful shutdown of the
	// main routine. Therefore, the ticker intervals should be configured so that the
	// time between ticker intervals is greater than the time it takes to complete the
	// longest running routine, to prevent routines from taking precedence over each
	// other.
	requestsTicker := time.NewTicker(c.conf.RequestInterval)
	reissuanceTicker := time.NewTicker(c.conf.ReissuanceInterval)
	log.Info().Dur("RequestInterval", c.conf.RequestInterval).Dur("ReissuanceInterval", c.conf.ReissuanceInterval).Str("store", c.certDir).Msg("cert-manager process started")

	for {
		// Wait for next tick or a stop signal
		select {
		case <-c.stop:
			log.Info().Msg("certificate manager received stop signal")
			return
		case <-requestsTicker.C:
			c.HandleCertificateRequests()
		case <-reissuanceTicker.C:
			c.HandleCertificateReissuance()
		}
	}

}

// HandleCertificateRequests performs one iteration through the certificate requests in
// the database and handles each sequentially, progressing them by modifying the status
// fields in the database. Note that this method logs errors instead of returning them
// to the caller.
func (c *CertificateManager) HandleCertificateRequests() {
	// Retrieve all certificate requests from the database
	var (
		nrequests int
		wg        sync.WaitGroup
		err       error
	)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	careqs := c.db.ListCertReqs(ctx)
	defer careqs.Release()
	log.Debug().Msg("cert-manager checking certificate request pipelines")

	for careqs.Next() {
		var req *models.CertificateRequest
		if req, err = careqs.CertReq(); err != nil {
			sentry.Error(nil).Err(err).Msg("could not parse certificate request from database")
			continue
		}

		logctx := sentry.With(nil).Str("id", req.Id).Str("common_name", req.CommonName)

		switch req.Status {
		case models.CertificateRequestState_READY_TO_SUBMIT:
			wg.Add(1)
			go func(req *models.CertificateRequest, logctx *sentry.Logger) {
				defer wg.Done()
				// Get the VASP from the certificate request
				var (
					vasp *pb.VASP
					err  error
				)
				if vasp, err = c.db.RetrieveVASP(ctx, req.Vasp); err != nil {
					logctx.Error().Str("vasp_id", req.Vasp).Msg("could not retrieve vasp for certificate request submission")
					return
				}

				// Verify that the VASP has not errored or been rejected
				if vasp.VerificationStatus < pb.VerificationState_REVIEWED || vasp.VerificationStatus > pb.VerificationState_VERIFIED {
					logctx.Info().Str("verification_status", vasp.VerificationStatus.String()).Msg("vasp is not ready for certificate request submission")
					if err = models.UpdateCertificateRequestStatus(req, models.CertificateRequestState_CR_REJECTED, "certificate request rejected", "automated"); err != nil {
						logctx.Error().Err(err).Msg("could not update certificate request status")
						return
					}
					if err = c.db.UpdateCertReq(ctx, req); err != nil {
						logctx.Error().Err(err).Msg("could not save updated certificate request")
						return
					}
				} else if err = c.submitCertificateRequest(req, vasp); err != nil {
					// If certificate submission requests fail we want immediate notification
					// so this is a CRITICAL severity that should alert us immediately.
					// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
					// this ensures that we issue a CRITICAL severity without stopping the server.
					sentry.Fatal(nil).Err(err).Msg("cert-manager could not submit certificate request")
				} else {
					logctx.Info().Msg("certificate request submitted")
				}
			}(req, logctx)
		case models.CertificateRequestState_PROCESSING:
			wg.Add(1)
			go func(req *models.CertificateRequest, logctx *sentry.Logger) {
				defer wg.Done()
				if err := c.checkCertificateRequest(req); err != nil {
					logctx.Error().Err(err).Msg("cert-manager could not process submitted certificate request")
				} else {
					logctx.Info().Msg("processing certificate request check complete")
				}
			}(req, logctx)
		}

		nrequests++
	}

	// Wait for all the certificate request processing to complete
	wg.Wait()

	// Check if there was an error while processing certificate requests
	if err = careqs.Error(); err != nil {
		sentry.Error(nil).Err(err).Msg("cert-manager could not retrieve certificate requests")
		return
	}

	// Conclude certificate handling successfully
	log.Debug().Int("requests", nrequests).Msg("cert-manager check complete")
}

func (c *CertificateManager) submitCertificateRequest(r *models.CertificateRequest, vasp *pb.VASP) (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Step 0: mark the VASP status as issuing certificates
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_ISSUING_CERTIFICATE, "issuing certificate", "automated"); err != nil {
		return err
	}
	if err = c.db.UpdateVASP(ctx, vasp); err != nil {
		return fmt.Errorf("could not update VASP status: %s", err)
	}

	// Step 1: find an authority with an available balance
	var authority int
	if authority, err = c.findCertAuthority(); err != nil {
		return err
	}

	// Step 2: get the password
	secretType := "password"
	pkcs12Password, err := c.secret.With(r.Id).GetLatestVersion(context.Background(), secretType)
	if err != nil {
		return fmt.Errorf("could not retrieve pkcs12password: %s", err)
	}

	// Allow multiple DNS names to be specified in addition to the common name
	// This will overwrite whatever is in the params ensuring the latest common name and
	// dns names are submitted to Sectigo if there were intermediate changes to the req.
	dnsNames := []string{r.CommonName}
	dnsNames = append(dnsNames, r.DnsNames...)
	models.UpdateCertificateRequestParams(r, sectigo.ParamDNSNames, strings.Join(dnsNames, "\n"))
	models.UpdateCertificateRequestParams(r, sectigo.ParamCommonName, r.CommonName)
	models.UpdateCertificateRequestParams(r, sectigo.ParamPassword, string(pkcs12Password))

	// Step 3: submit the certificate
	var (
		rep    *sectigo.BatchResponse
		params map[string]string
	)

	// Construct the required parameters for the Sectigo request.
	profile := c.certs.Profile()
	batchName := fmt.Sprintf("%s-certreq-%s", c.conf.DirectoryID, r.Id)

	if params, err = models.GetCertificateRequestParams(r, profile); err != nil {
		return fmt.Errorf("could not retrieve certificate request parameters for profile %q: %s", profile, err)
	}

	// Execute the certificate request to Sectigo.
	if rep, err = c.certs.CreateSingleCertBatch(authority, batchName, params); err != nil {
		// Although the error may be logged again by the calling function, log the error
		// here as well to provide debugging information about why the Sectigo request failed.
		dict := sentry.Dict()
		for key, value := range params {
			// NOTE: Do not log any passwords or secrets!
			if key == "pkcs12Password" {
				value = strings.Repeat("*", len(value))
			}
			dict.Str(key, value)
		}
		sentry.Error(nil).Err(err).
			Int("authority", authority).
			Str("batch_name", batchName).
			Dict("params", dict).
			Str("profile", profile).
			Msg("create single cert batch failed")
		return fmt.Errorf("could not create single certificate batch: %s", err)
	}

	// Step 4: update the certificate request with the batch details
	r.AuthorityId = int64(authority)
	r.BatchId = int64(rep.BatchID)
	r.BatchName = rep.BatchName
	r.BatchStatus = rep.Status
	r.OrderNumber = int64(rep.OrderNumber)
	r.CreationDate = rep.CreationDate
	r.Profile = rep.Profile
	r.RejectReason = rep.RejectReason

	// Mark the certificate request as processing so downstream status checks occur
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_PROCESSING, "certificate submitted", "automated"); err != nil {
		return fmt.Errorf("could not update certificate request status: %s", err)
	}
	if err = c.db.UpdateCertReq(ctx, r); err != nil {
		return fmt.Errorf("could not update certificate with batch details: %s", err)
	}

	return nil
}

func (c *CertificateManager) checkCertificateRequest(r *models.CertificateRequest) (err error) {
	if r.BatchId == 0 {
		return errors.New("missing batch ID - cannot retrieve status")
	}

	// Step 1: refresh batch info from the sectigo service
	var info *sectigo.BatchResponse
	if info, err = c.certs.BatchDetail(int(r.BatchId)); err != nil {
		return fmt.Errorf("could not fetch batch info for id %d: %s", r.BatchId, err)
	}

	// Step 1b: update certificate request with fetched info
	r.BatchStatus = info.Status
	r.RejectReason = info.RejectReason

	// Step 1c: check if the batch is in an unhandled state, and if so, refresh batch status
	if info.Status == sectigo.BatchStatusCollected || info.Status == "" {
		sentry.Warn(nil).Int64("batch_id", r.BatchId).Str("batch_status", info.Status).Msg("unknown batch info status, refreshing batch status directly")
		if r.BatchStatus, err = c.certs.BatchStatus(int(r.BatchId)); err != nil {
			return fmt.Errorf("could not fetch batch status for id %d: %s", r.BatchId, err)
		}
	}

	// Step 2: get the processing info for the batch
	var proc *sectigo.ProcessingInfoResponse
	if proc, err = c.certs.ProcessingInfo(int(r.BatchId)); err != nil {
		return fmt.Errorf("could not fetch batch processing info for id %d: %s", r.BatchId, err)
	}

	log.Info().
		Str("status", r.BatchStatus).
		Str("reject", r.RejectReason).
		Int("active", proc.Active).
		Int("failed", proc.Failed).
		Int("success", proc.Success).
		Msg("batch processing status")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Step 3: check active - if there is still an active batch then delay
	if proc.Active > 0 {
		if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_PROCESSING, "awaiting batch processing", "automated"); err != nil {
			return fmt.Errorf("could not update certificate request status: %s", err)
		}
		if err = c.db.UpdateCertReq(ctx, r); err != nil {
			return fmt.Errorf("could not save updated cert request: %s", err)
		}
		return nil
	}

	// Step 4: check failures -- determine if certificate request has been rejected
	if proc.Failed > 0 {
		logctx := sentry.With(nil).
			Int("batch_id", int(r.BatchId)).
			Int("failed", proc.Failed).
			Int("success", proc.Success).
			Str("status", r.BatchStatus).
			Str("name", r.BatchName)

		if proc.Success > 0 || r.BatchStatus == sectigo.BatchStatusReadyForDownload {
			// This may mean that some certificates can be downloaded, so just log
			// errors and continue with download processing
			logctx.Warn().Msg("certificate request mixed success/failure")
		} else {
			// In this case there were no successes, so set certificate request status accordingly
			// and do not continue processing the certificate request
			if r.RejectReason != "" || r.BatchStatus == sectigo.BatchStatusRejected {
				// Assume the certificate was rejected
				if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_CR_REJECTED, "certificate request rejected", "automated"); err != nil {
					return fmt.Errorf("could not update certificate request status: %s", err)
				}
				logctx.Warn().Msg("certificate request rejected")
			} else {
				// Assume the certificate errored and wasn't rejected
				if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_CR_ERRORED, "certificate request errored", "automated"); err != nil {
					return fmt.Errorf("could not update certificate request status: %s", err)
				}
				logctx.Warn().Msg("certificate request errored")
			}

			if err = c.db.UpdateCertReq(ctx, r); err != nil {
				return fmt.Errorf("could not save updated cert request: %s", err)
			}
			return nil
		}
	}

	// Step 5: Check to make sure we can download certificates
	if proc.Success == 0 || r.BatchStatus != sectigo.BatchStatusReadyForDownload {
		// We should not be in this state, it should have been handled in Step 1c or 4
		// so this is a developer error on our part, or a change in the Sectigo API
		// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
		// this ensures that we issue a CRITICAL severity without stopping the server.
		sentry.Fatal(nil).Int64("batch_id", r.BatchId).Int("success", proc.Success).Str("batch_status", r.BatchStatus).Msg("unhandled sectigo state")
		if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_PROCESSING, "unhandled sectigo state", "automated"); err != nil {
			return fmt.Errorf("could not update certificate request status: %s", err)
		}
		if err = c.db.UpdateCertReq(ctx, r); err != nil {
			return fmt.Errorf("could not save updated cert request: %s", err)
		}
		return nil
	}

	// Step 6: Mark the status as ready for download!
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_DOWNLOADING, "certificate ready for download", "automated"); err != nil {
		return fmt.Errorf("could not update certificate request status: %s", err)
	}
	if err = c.db.UpdateCertReq(ctx, r); err != nil {
		return fmt.Errorf("could not save updated cert request: %s", err)
	}

	// Fetch the VASP from the certificate request
	var vasp *pb.VASP
	if vasp, err = c.db.RetrieveVASP(ctx, r.Vasp); err != nil {
		return fmt.Errorf("could not retrieve vasp: %s", err)
	}

	// Make sure the VASP has not errored or been rejected before downloading certificates
	if vasp.VerificationStatus != pb.VerificationState_ISSUING_CERTIFICATE {
		sentry.Error(c).Str("verification_status", vasp.VerificationStatus.String()).Msg("VASP is not in the ISSUING_CERTIFICATE state, cannot download certificates")
		if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_CR_REJECTED, "rejecting certificate request", "automated"); err != nil {
			return fmt.Errorf("could not update certificate request status: %s", err)
		}
		if err = c.db.UpdateCertReq(ctx, r); err != nil {
			return fmt.Errorf("could not save updated cert request: %s", err)
		}
		return nil
	}

	// Send off downloader go routine to fetch the certs and notify the user
	c.downloadCertificateRequest(r, vasp)
	return nil
}

// finds the first authority with an available balance greater than 0.
func (c *CertificateManager) findCertAuthority() (id int, err error) {
	var authorities []*sectigo.AuthorityResponse
	if authorities, err = c.certs.UserAuthorities(); err != nil {
		return 0, fmt.Errorf("could not fetch user authorities: %s", err)
	}

	for _, authority := range authorities {
		var balance int
		if balance, err = c.certs.AuthorityAvailableBalance(authority.ID); err != nil {
			sentry.Error(c).Err(err).Int("authority", authority.ID).Msg("could not fetch authority balance")
		}
		if balance > 0 {
			return authority.ID, nil
		}
	}

	return 0, fmt.Errorf("could not find authority with available balance out of %d available authorities", len(authorities))
}

// a go routine that downloads the certificate in the background, then sends the certs
// as an attachment to the technical contact if available.
func (c *CertificateManager) downloadCertificateRequest(r *models.CertificateRequest, vasp *pb.VASP) {
	var (
		err        error
		path       string
		payload    []byte
		secretType string
	)

	// Download the certificates as a zip file to the cert storage directory
	if path, err = c.certs.Download(int(r.BatchId), c.certDir); err != nil {
		sentry.Error(c).Err(err).Int("batch", int(r.BatchId)).Msg("could not download certificates")
		return
	}

	// Store zipped cert in Google Secrets using certRequestId
	sctx := context.Background()
	secretType = "cert"
	if err = c.secret.With(r.Id).CreateSecret(sctx, secretType); err != nil {
		sentry.Error(c).Err(err).Msg("could not create cert secret")
		return
	}
	if payload, err = os.ReadFile(path); err != nil {
		sentry.Error(nil).Err(err).Msg("could not read in cert payload")
		return
	}
	if err = c.secret.With(r.Id).AddSecretVersion(sctx, secretType, payload); err != nil {
		sentry.Error(nil).Err(err).Msg("could not add secret version for cert payload")
		return
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Mark as downloaded.
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_DOWNLOADED, "certificate downloaded", "automated"); err != nil {
		sentry.Error(nil).Err(err).Msg("could not update certificate request status")
		return
	}
	if err = c.db.UpdateCertReq(ctx, r); err != nil {
		sentry.Error(nil).Err(err).Msg("could not save updated cert request")
		return
	}

	// Delete the temporary file
	defer os.Remove(path)

	log.Info().Str(path, path).Msg("certificates written to secret manager")

	// Retrieve the latest secret version for the password
	secretType = "password"
	pkcs12password, err := c.secret.With(r.Id).GetLatestVersion(sctx, secretType)
	if err != nil {
		sentry.Error(nil).Err(err).Msg("could not retrieve password from secret manager to extract public key")
		return
	}

	if vasp.IdentityCertificate, err = extractCertificate(path, string(pkcs12password)); err != nil {
		sentry.Error(nil).Err(err).Msg("could not extract certificate")
		return
	}

	// Create the certificate record
	var cert *models.Certificate
	if cert, err = models.NewCertificate(vasp, r, vasp.IdentityCertificate); err != nil {
		sentry.Error(nil).Err(err).Msg("could not create certificate record")
		return
	}

	// Update the certificate record
	if err = c.db.UpdateCert(ctx, cert); err != nil {
		sentry.Error(nil).Err(err).Msg("could not update certificate record")
		return
	}

	// Add the certificate ID to the request and VASP records
	r.Certificate = cert.Id
	if err = models.AppendCertID(vasp, cert.Id); err != nil {
		sentry.Error(nil).Err(err).Msg("could not append certificate ID to VASP")
		return
	}

	// Update the VASP status as verified/certificate issued
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_VERIFIED, "certificate issued", "automated"); err != nil {
		sentry.Error(nil).Err(err).Msg("could not update VASP verification status")
		return
	}
	if err = c.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(nil).Err(err).Msg("could not update VASP status as verified")
		return
	}

	var deliveryErr error
	if r.Webhook != "" {
		// Deliver the certificates payload using the configured webhook.
		if deliveryErr = c.deliverCertificatePayload(ctx, r, payload); deliveryErr != nil {
			log.Error().Err(deliveryErr).Str("webhook", r.Webhook).Str("vasp", vasp.Id).Msg("error delivering certificate via webhook")
		}
	}

	// If the user has not specifically turned off email delivery or if webhook
	// delivery failed, send the certificates via email.
	if !r.NoEmailDelivery || deliveryErr != nil {
		if _, err = c.email.SendDeliverCertificates(vasp, path); err != nil {
			// If there is an error delivering emails, return here so we don't mark as completed
			sentry.Error(nil).Err(err).Msg("could not deliver certificates to technical contact")
			return
		}
	}

	if err = c.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(nil).Err(err).Msg("could not update VASP email logs")
		return
	}

	// Mark certificate request as complete.
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_COMPLETED, "certificate request complete", "automated"); err != nil {
		sentry.Error(nil).Err(err).Msg("could not update certificate request status")
		return
	}
	r.Status = models.CertificateRequestState_COMPLETED
	if err = c.db.UpdateCertReq(ctx, r); err != nil {
		sentry.Error(nil).Err(err).Msg("could not save updated cert request")
		return
	}

	log.Info().
		Str("serial_number", hex.EncodeToString(vasp.IdentityCertificate.SerialNumber)).
		Msg("certificates extracted and delivered")
}

func extractCertificate(path, pkcs12password string) (pub *pb.Certificate, err error) {
	var archive *trust.Serializer
	if archive, err = trust.NewSerializer(true, pkcs12password, trust.CompressionZIP); err != nil {
		return nil, err
	}

	var provider *trust.Provider
	if provider, err = archive.ReadFile(path); err != nil {
		return nil, err
	}

	var cert *x509.Certificate
	if cert, err = provider.GetLeafCertificate(); err != nil {
		return nil, err
	}

	pub = &pb.Certificate{
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
		return nil, fmt.Errorf("could not PEM encode certificate: %s", err)
	}
	pub.Data = buf.Bytes()

	// Write the entire provider chain into the directory service data store
	if archive, err = trust.NewSerializer(false, "", trust.CompressionGZIP); err != nil {
		return nil, err
	}

	// Ensure only the public keys are written to the directory service
	if pub.Chain, err = archive.Compress(provider.Public()); err != nil {
		return nil, err
	}

	return pub, nil
}

// get the configured cert storage directory or return a temporary directory
func (c *CertificateManager) getCertStorage() (path string, err error) {
	if c.conf.Storage != "" {
		var stat os.FileInfo
		if stat, err = os.Stat(c.conf.Storage); err != nil {
			if os.IsNotExist(err) {
				// Create the directory if it does not exist and return
				if err = os.MkdirAll(c.conf.Storage, 0755); err != nil {
					return "", fmt.Errorf("could not create cert storage directory: %s", err)
				}
				return c.conf.Storage, nil
			}

			// Other permissions error, cannot access cert storage
			return "", err
		}

		if !stat.IsDir() {
			return "", errors.New("not a directory")
		}
		return c.conf.Storage, nil
	}

	// Create a temporary directory
	if path, err = os.MkdirTemp("", "gds_certs"); err != nil {
		return "", err
	}
	log.Warn().Str("certs", path).Msg("using a temporary directory for cert downloads")
	return path, err
}

// HandleCertificateReissuance iterates through each VASP in the database and checks
// if their identity certificate will be expiring soon, sending a reminder email at
// the 30 and 7 day checkpoints if so, and reissuing the identity certificate 10 days before
// expiration.
func (c *CertificateManager) HandleCertificateReissuance() {
	var (
		err            error
		expirationDate time.Time
	)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Iterate through the VASPs in the database.
	vasps := c.db.ListVASPs(ctx)
	defer vasps.Release()
vaspsLoop:
	for vasps.Next() {
		var vasp *pb.VASP
		if vasp, err = vasps.VASP(); err != nil {
			sentry.Error(nil).Err(err).Msg("could not parse vasp record from database")
			continue vaspsLoop
		}

		// Skip the current VASP if it is not verified.
		if vasp.GetVerificationStatus() != pb.VerificationState_VERIFIED {
			continue vaspsLoop
		}

		// Make sure the VASP has an identity certificate to avoid a panic.
		identityCert := vasp.IdentityCertificate
		if identityCert == nil {
			sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("vasp is verified but does not have an identity certificate")
			continue vaspsLoop
		}

		// Calculate the number of days before the VASP's certificate expires.
		if expirationDate, err = time.Parse(time.RFC3339, identityCert.NotAfter); err != nil {
			sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not parse %s's cert reissuance date")
			continue vaspsLoop
		}

		// Compute VASP signature to check if it has been modified
		var sig []byte
		if sig, err = models.VASPSignature(vasp); err != nil {
			sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not compute signature for VASP")
			continue vaspsLoop
		}

		// Calculate the reissuance date, 10 days before expiration
		reissuanceDate := expirationDate.Add(-time.Hour * 240)
		timeBeforeExpiration := time.Until(expirationDate)

		// NOTE: This computation returns fractional days rather than rounding up or down to the nearest day.
		switch daysBeforeExpiration := timeBeforeExpiration.Hours() / 24; {

		// If a certificate has expired, update the certificate record.
		case daysBeforeExpiration <= 0:
			var cert *models.Certificate
			if cert, err = c.db.RetrieveCert(ctx, models.GetCertID(identityCert)); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not retrieve expired certificate record")
				continue vaspsLoop
			}
			cert.Status = models.CertificateState_EXPIRED
			if err = c.db.UpdateCert(ctx, cert); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not update expired certificate record status")
			}

		// Seven days before expiration, send a cert reissuance reminder to VASP.
		// NOTE: the SendContactReissuanceReminder will not send emails more than
		// once to a contact.
		case daysBeforeExpiration <= 7:
			if err = c.email.SendContactReissuanceReminder(vasp, 7, reissuanceDate); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("error sending seven day reissuance reminder")
				continue vaspsLoop
			}

		// Ten days before expiration, reissue the VASP's identity certificate, send the email with the created pkcs12
		// password and send the whisper link, as well as notifying the TRISA admin that reissuance has started.
		case daysBeforeExpiration <= 10:
			// Check if the reissuance process has already started for this VASP
			if started, err := c.reissuanceInProgress(vasp); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not check vasp reissuance status")
				continue vaspsLoop
			} else if started {
				log.Info().Str("vasp_id", vasp.Id).Msg("vasp reissuance is already in progress")
				continue vaspsLoop
			}

			// Start the reissuance process for this VASP by creating a new certificate request
			if err = c.reissueIdentityCertificates(vasp); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not start reissuance process")
				continue vaspsLoop
			}
			if _, err = c.email.SendReissuanceAdminNotification(vasp, 10, reissuanceDate); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("error sending admin reissuance notification")
				// The VASP and certreq records have already been updated in
				// reissueIdentityCertificates, so we can continue to the next VASP if
				// we failed to send the email.
				continue vaspsLoop
			}

		// Thirty days before expiration, send the reissuance reminder to the VASP and the TRISA admin.
		case daysBeforeExpiration <= 30:
			// If the reminder fails do not stop processing and attempt to send reminder to admins
			if err = c.email.SendContactReissuanceReminder(vasp, 30, reissuanceDate); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("error sending thirty day reissuance reminder")
			}

			if _, err = c.email.SendExpiresAdminNotification(vasp, 30, reissuanceDate); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("error sending admin reissuance reminder")
				continue vaspsLoop
			}
		}

		// Perform VASP hash check to determine if the VASP has been updated, and if it
		// has, save the updates back to the database (otherwise do not save so that the
		// last modified timestamp is not updated every day).
		var updated []byte
		if updated, err = models.VASPSignature(vasp); err != nil {
			sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not compute signature for VASP")
			continue vaspsLoop
		}

		if !bytes.Equal(sig, updated) {
			// We need to update the vasp record in the database so that the email logs are preserved.
			if err = c.db.UpdateVASP(ctx, vasp); err != nil {
				sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("error updating the VASP record in the database")
				continue vaspsLoop
			}
		}
	}

	if err := vasps.Error(); err != nil {
		sentry.Error(nil).Err(err).Msg("could not iterate through database")
	}
}

// Helper to check if the reissuance process has already started for a VASP.
func (c *CertificateManager) reissuanceInProgress(vasp *pb.VASP) (_ bool, err error) {
	// Get the latest certificate request ID for the VASP.
	var certReqID string
	if certReqID, err = models.GetLatestCertReqID(vasp); err != nil {
		return false, err
	}

	// If there are no existing certificate requests, then the reissuance process has not started.
	if certReqID == "" {
		return false, nil
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Get the certificate request from the database.
	var certReq *models.CertificateRequest
	if certReq, err = c.db.RetrieveCertReq(ctx, certReqID); err != nil {
		return false, err
	}

	return certReq.Status >= models.CertificateRequestState_READY_TO_SUBMIT && certReq.Status < models.CertificateRequestState_COMPLETED, nil
}

// Helper function for HandleCertificateReissuance that reissues identity certificates
// for the given vasp.
func (c *CertificateManager) reissueIdentityCertificates(vasp *pb.VASP) (err error) {
	var (
		certreq     *models.CertificateRequest
		whisperLink string
	)

	// Create a new certificate request.
	if certreq, err = models.NewCertificateRequest(vasp); err != nil {
		return fmt.Errorf("error creating certificate request for vasp %s: %w", vasp.Id, err)
	}

	// Update the cert req to be ready for submission.
	if err = models.UpdateCertificateRequestStatus(
		certreq,
		models.CertificateRequestState_READY_TO_SUBMIT,
		"automated certificate reissuance",
		"automated",
	); err != nil {
		return fmt.Errorf("error updating certificate request for vasp %s: %w", vasp.Id, err)
	}

	// Generate a new PKCS12 password
	secretType := "password"
	pkcs12password := secrets.CreateToken(16)
	whisperPasswordTemplate := "Below is the PKCS12 password which you must use to decrypt your new certificates:\n\n%s\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new secret using the secret manager.
	if err = c.secret.With(certreq.Id).CreateSecret(ctx, secretType); err != nil {
		return fmt.Errorf("error creating password secret for vasp %s: %w", vasp.Id, err)
	}
	if err = c.secret.With(certreq.Id).AddSecretVersion(ctx, secretType, []byte(pkcs12password)); err != nil {
		return fmt.Errorf("error creating password version for vasp %s: %w", vasp.Id, err)
	}

	var deliveryErr error
	if certreq.Webhook != "" {
		// Deliver the PKCS12 password using the configured webhook.
		if deliveryErr = c.deliverCertificatePassword(ctx, certreq, pkcs12password); deliveryErr != nil {
			log.Error().Err(deliveryErr).Str("webhook", certreq.Webhook).Str("vasp", vasp.Id).Msg("error delivering pkcs12 password via webhook")
		}
	}

	// If the user has not specifically turned off email delivery or if there was an
	// error in webhook delivery, send the pkcs12 password in a whisper via email.
	if !certreq.NoEmailDelivery || deliveryErr != nil {
		if whisperLink, err = whisper.CreateSecretLink(fmt.Sprintf(whisperPasswordTemplate, pkcs12password), "", 3, time.Now().AddDate(0, 0, 7)); err != nil {
			return fmt.Errorf("error creating whisper link for vasp %s: %w", vasp.Id, err)
		}

		if _, err = c.email.SendReissuanceStarted(vasp, whisperLink); err != nil {
			return fmt.Errorf("error sending reissuance started email for vasp %s: %w", vasp.Id, err)
		}
	}

	// Update the certificate request in the datastore.
	if err = c.db.UpdateCertReq(ctx, certreq); err != nil {
		return fmt.Errorf("error updating certificate request for vasp %s: %w", vasp.Id, err)
	}

	// Save the certificate request on the VASP.
	if err = models.AppendCertReqID(vasp, certreq.Id); err != nil {
		return fmt.Errorf("error appending certificate request to vasp %s: %w", vasp.Id, err)
	}

	// Update the VASP information in the datastore.
	if err = c.db.UpdateVASP(ctx, vasp); err != nil {
		return fmt.Errorf("error updating vasp %s in the certman store: %w", vasp.Id, err)
	}
	return nil
}

// Attempt to deliver a certificate password by webhook using the configured backoff
// strategy or return an error if the maximum number of retries was exceeded.
func (c *CertificateManager) deliverCertificatePassword(ctx context.Context, certreq *models.CertificateRequest, password string) (err error) {
	// Create the client to deliver the password.
	var client courier.CourierClient
	if client, err = courier.New(certreq.Webhook); err != nil {
		return fmt.Errorf("could not create courier client for pkcs12 password delivery: %s", err)
	}

	req := &courier.StorePasswordRequest{
		ID:       certreq.Id,
		Password: password,
	}

	// Attempt to deliver the password using the configured backoff strategy.
	wait := c.conf.DeliveryBackoff.Ticker()
	defer wait.Stop()

	// Wait for the context to be cancelled or the password to be delivered.
	var retries int
	for retries = 0; retries < c.conf.DeliveryBackoff.MaxRetries+1; retries++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wait.C:
			if err = client.StoreCertificatePassword(ctx, req); err != nil {
				log.Warn().Err(err).Int("attempts", retries+1).Str("webhook", certreq.Webhook).Msg("could not deliver certificate password, retrying")
				continue
			}
			return nil
		}
	}

	// Exceeded the maximum number of retries
	return fmt.Errorf("could not deliver certificate password after %d retries", retries)
}

// Deliver a certificate payload by webhook using the configured backoff strategy or
// return an error if the maximum number of retries was exceeded.
func (c *CertificateManager) deliverCertificatePayload(ctx context.Context, certreq *models.CertificateRequest, payload []byte) (err error) {
	// Create the client to deliver the password.
	var client courier.CourierClient
	if client, err = courier.New(certreq.Webhook); err != nil {
		return fmt.Errorf("could not create courier client for certificate delivery: %s", err)
	}

	req := &courier.StoreCertificateRequest{
		ID:                certreq.Id,
		Base64Certificate: base64.StdEncoding.EncodeToString(payload),
	}

	// Attempt to deliver the password using the configured backoff strategy.
	wait := c.conf.DeliveryBackoff.Ticker()
	defer wait.Stop()

	// Wait for the context to be cancelled or the password to be delivered.
	var retries int
	for retries = 0; retries < c.conf.DeliveryBackoff.MaxRetries+1; retries++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wait.C:
			if err = client.StoreCertificate(ctx, req); err != nil {
				log.Warn().Err(err).Int("attempts", retries+1).Str("webhook", certreq.Webhook).Msg("could not deliver encrypted certificate, retrying")
				continue
			}
			return nil
		}
	}

	// Exceeded the maximum number of retries
	return fmt.Errorf("could not deliver encrypted certificate after %d retries", retries)
}
