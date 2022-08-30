package certman

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
	"github.com/trisacrypto/directory/pkg/utils/whisper"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func New(conf config.CertManConfig, db store.Store, secret *secrets.SecretManager, email *emails.EmailManager, directoryID string) (cm *CertificateManager, err error) {
	cm = &CertificateManager{
		conf:        conf,
		db:          db,
		secret:      secret,
		email:       email,
		directoryID: directoryID,
	}

	if cm.directoryID == "" {
		return nil, errors.New("directory ID is required for cert manager")
	}

	if cm.email == nil {
		return nil, errors.New("email manager is required for cert manager")
	}

	if cm.secret == nil {
		return nil, errors.New("secret manager is required for cert manager")
	}

	if conf.Sectigo.Testing {
		if err = mock.Start(conf.Sectigo.Profile); err != nil {
			return nil, err
		}
	}

	if cm.certs, err = sectigo.New(conf.Sectigo); err != nil {
		return nil, err
	}

	if _, err = cm.getCertStorage(); err != nil {
		return nil, err
	}

	return cm, nil
}

// CertificateManager is a struct with a go routine that periodically checks on the
// status of certificate requests and moves them through the request pipeline. This is
// separated from the parent GDS to allow for isolated testing.
type CertificateManager struct {
	conf        config.CertManConfig
	db          store.Store
	secret      *secrets.SecretManager
	certs       *sectigo.Sectigo
	email       *emails.EmailManager
	directoryID string
}

// Run starts the CertManager go routine with the given channel and returns a wait
// group, allowing the caller to control synchronization from the outside.
func (c *CertificateManager) Run(stop chan struct{}) (wg *sync.WaitGroup) {
	if stop == nil {
		stop = make(chan struct{})
	}
	wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.CertManager(stop)
	}()
	return wg
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
func (c *CertificateManager) CertManager(stop <-chan struct{}) {
	// Check certificate download directory
	certDir, err := c.getCertStorage()
	if err != nil {
		log.Fatal().Err(err).Msg("cert-manager cannot access certificate storage")
	}

	// Ticker is created in the go routine to prevent backpressure if the cert manager
	// process takes longer than the specified ticker interval.
	requestsTicker := time.NewTicker(c.conf.RequestInterval)
	reissuanceTicker := time.NewTicker(c.conf.ReissuenceInterval)
	log.Info().Dur("RequestInterval", c.conf.RequestInterval).Dur("ReissuenceInterval", c.conf.ReissuenceInterval).Str("store", certDir).Msg("cert-manager process started")

	for {
		// Wait for next tick or a stop signal
		select {
		case <-stop:
			log.Info().Msg("certificate manager received stop signal")
			return
		case <-requestsTicker.C:
			c.HandleCertificateRequests(certDir)
		case <-reissuanceTicker.C:
			c.HandleCertificateReissuance()
		}
	}

}

// HandleCertificateRequests performs one iteration through the certificate requests in
// the database and handles each sequentially, progressing them by modifying the status
// fields in the database. Note that this method logs errors instead of returning them
// to the caller.
func (c *CertificateManager) HandleCertificateRequests(certDir string) {
	// Retrieve all certificate requests from the database
	var (
		nrequests int
		wg        sync.WaitGroup
		err       error
	)

	careqs := c.db.ListCertReqs()
	defer careqs.Release()
	log.Debug().Msg("cert-manager checking certificate request pipelines")

	for careqs.Next() {
		var req *models.CertificateRequest
		if req, err = careqs.CertReq(); err != nil {
			log.Error().Err(err).Msg("could not parse certificate request from database")
			continue
		}

		logctx := log.With().Str("id", req.Id).Str("common_name", req.CommonName).Logger()

		switch req.Status {
		case models.CertificateRequestState_READY_TO_SUBMIT:
			wg.Add(1)
			go func(req *models.CertificateRequest, logctx zerolog.Logger) {
				defer wg.Done()
				// Get the VASP from the certificate request
				var (
					vasp *pb.VASP
					err  error
				)
				if vasp, err = c.db.RetrieveVASP(req.Vasp); err != nil {
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
					if err = c.db.UpdateCertReq(req); err != nil {
						logctx.Error().Err(err).Msg("could not save updated certificate request")
						return
					}
				} else if err = c.submitCertificateRequest(req, vasp); err != nil {
					// If certificate submission requests fail we want immediate notification
					// so this is a CRITICAL severity that should alert us immediately.
					// NOTE: using WithLevel and Fatal does not Exit the program like log.Fatal()
					// this ensures that we issue a CRITICAL severity without stopping the server.
					log.WithLevel(zerolog.FatalLevel).Err(err).Msg("cert-manager could not submit certificate request")
				} else {
					logctx.Info().Msg("certificate request submitted")
				}
			}(req, logctx)
		case models.CertificateRequestState_PROCESSING:
			wg.Add(1)
			go func(req *models.CertificateRequest, logctx zerolog.Logger) {
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
		log.Error().Err(err).Msg("cert-manager could not retrieve certificate requests")
		return
	}

	// Conclude certificate handling successfully
	log.Debug().Int("requests", nrequests).Msg("cert-manager check complete")
}

func (c *CertificateManager) submitCertificateRequest(r *models.CertificateRequest, vasp *pb.VASP) (err error) {
	// Step 0: mark the VASP status as issuing certificates
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_ISSUING_CERTIFICATE, "issuing certificate", "automated"); err != nil {
		return err
	}
	if err = c.db.UpdateVASP(vasp); err != nil {
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

	profile := c.certs.Profile()
	var params map[string]string
	if profile == sectigo.ProfileCipherTraceEndEntityCertificate || profile == sectigo.ProfileIDCipherTraceEndEntityCertificate {
		params = r.Params
		if params == nil {
			log.Error().Str("vasp", vasp.Id).Str("certreq", r.Id).Msg("certificate request params are nil")
			return errors.New("no params are available on the certificate request")
		}
	} else {
		params = make(map[string]string)
	}

	params["commonName"] = r.CommonName
	params["dNSName"] = r.CommonName
	params["pkcs12Password"] = string(pkcs12Password)

	// Step 3: submit the certificate
	var rep *sectigo.BatchResponse
	batchName := fmt.Sprintf("%s-certreq-%s)", c.directoryID, r.Id)
	if rep, err = c.certs.CreateSingleCertBatch(authority, batchName, params); err != nil {
		// Although the error may be logged again by the calling function, log the error
		// here as well to provide debugging information about why the Sectigo request failed.
		dict := zerolog.Dict()
		for key, value := range params {
			// NOTE: Do not log any passwords or secrets!
			if key == "pkcs12Password" {
				value = strings.Repeat("*", len(value))
			}
			dict.Str(key, value)
		}
		log.Error().Err(err).
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
	if err = c.db.UpdateCertReq(r); err != nil {
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
		log.Warn().Int64("batch_id", r.BatchId).Str("batch_status", info.Status).Msg("unknown batch info status, refreshing batch status directly")
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

	// Step 3: check active - if there is still an active batch then delay
	if proc.Active > 0 {
		if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_PROCESSING, "awaiting batch processing", "automated"); err != nil {
			return fmt.Errorf("could not update certificate request status: %s", err)
		}
		if err = c.db.UpdateCertReq(r); err != nil {
			return fmt.Errorf("could not save updated cert request: %s", err)
		}
		return nil
	}

	// Step 4: check failures -- determine if certificate request has been rejected
	if proc.Failed > 0 {
		logctx := log.With().
			Int("batch_id", int(r.BatchId)).
			Int("failed", proc.Failed).
			Int("success", proc.Success).
			Str("status", r.BatchStatus).
			Str("name", r.BatchName).
			Logger()

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

			if err = c.db.UpdateCertReq(r); err != nil {
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
		log.WithLevel(zerolog.FatalLevel).Int64("batch_id", r.BatchId).Int("success", proc.Success).Str("batch_status", r.BatchStatus).Msg("unhandled sectigo state")
		if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_PROCESSING, "unhandled sectigo state", "automated"); err != nil {
			return fmt.Errorf("could not update certificate request status: %s", err)
		}
		if err = c.db.UpdateCertReq(r); err != nil {
			return fmt.Errorf("could not save updated cert request: %s", err)
		}
		return nil
	}

	// Step 6: Mark the status as ready for download!
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_DOWNLOADING, "certificate ready for download", "automated"); err != nil {
		return fmt.Errorf("could not update certificate request status: %s", err)
	}
	if err = c.db.UpdateCertReq(r); err != nil {
		return fmt.Errorf("could not save updated cert request: %s", err)
	}

	// Fetch the VASP from the certificate request
	var vasp *pb.VASP
	if vasp, err = c.db.RetrieveVASP(r.Vasp); err != nil {
		return fmt.Errorf("could not retrieve vasp: %s", err)
	}

	// Make sure the VASP has not errored or been rejected before downloading certificates
	if vasp.VerificationStatus != pb.VerificationState_ISSUING_CERTIFICATE {
		log.Error().Str("verification_status", vasp.VerificationStatus.String()).Msg("VASP is not in the ISSUING_CERTIFICATE state, cannot download certificates")
		if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_CR_REJECTED, "rejecting certificate request", "automated"); err != nil {
			return fmt.Errorf("could not update certificate request status: %s", err)
		}
		if err = c.db.UpdateCertReq(r); err != nil {
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
			log.Error().Err(err).Int("authority", authority.ID).Msg("could not fetch authority balance")
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
		certDir    string
		payload    []byte
		secretType string
	)

	// Get the cert storage directory to download certs to
	if certDir, err = c.getCertStorage(); err != nil {
		log.Error().Err(err).Msg("could not find cert storage directory")
		return
	}

	// Download the certificates as a zip file to the cert storage directory
	if path, err = c.certs.Download(int(r.BatchId), certDir); err != nil {
		log.Error().Err(err).Int("batch", int(r.BatchId)).Msg("could not download certificates")
		return
	}

	// Store zipped cert in Google Secrets using certRequestId
	sctx := context.Background()
	secretType = "cert"
	if err = c.secret.With(r.Id).CreateSecret(sctx, secretType); err != nil {
		log.Error().Err(err).Msg("could not create cert secret")
		return
	}
	if payload, err = ioutil.ReadFile(path); err != nil {
		log.Error().Err(err).Msg("could not read in cert payload")
		return
	}
	if err = c.secret.With(r.Id).AddSecretVersion(sctx, secretType, payload); err != nil {
		log.Error().Err(err).Msg("could not add secret version for cert payload")
		return
	}

	// Mark as downloaded.
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_DOWNLOADED, "certificate downloaded", "automated"); err != nil {
		log.Error().Err(err).Msg("could not update certificate request status")
		return
	}
	if err = c.db.UpdateCertReq(r); err != nil {
		log.Error().Err(err).Msg("could not save updated cert request")
		return
	}

	// Delete the temporary file
	defer os.Remove(path)

	log.Info().Str(path, path).Msg("certificates written to secret manager")

	// Retrieve the latest secret version for the password
	secretType = "password"
	pkcs12password, err := c.secret.With(r.Id).GetLatestVersion(sctx, secretType)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve password from secret manager to extract public key")
		return
	}

	if vasp.IdentityCertificate, err = extractCertificate(path, string(pkcs12password)); err != nil {
		log.Error().Err(err).Msg("could not extract certificate")
		return
	}

	// Create the certificate record
	var cert *models.Certificate
	if cert, err = models.NewCertificate(vasp, r, vasp.IdentityCertificate); err != nil {
		log.Error().Err(err).Msg("could not create certificate record")
		return
	}

	// Update the certificate record
	if err = c.db.UpdateCert(cert); err != nil {
		log.Error().Err(err).Msg("could not update certificate record")
		return
	}

	// Add the certificate ID to the request and VASP records
	r.Certificate = cert.Id
	if err = models.AppendCertID(vasp, cert.Id); err != nil {
		log.Error().Err(err).Msg("could not append certificate ID to VASP")
		return
	}

	// Update the VASP status as verified/certificate issued
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_VERIFIED, "certificate issued", "automated"); err != nil {
		log.Error().Err(err).Msg("could not update VASP verification status")
		return
	}
	if err = c.db.UpdateVASP(vasp); err != nil {
		log.Error().Err(err).Msg("could not update VASP status as verified")
		return
	}

	// Email the certificates to the technical contacts
	if _, err = c.email.SendDeliverCertificates(vasp, path); err != nil {
		// If there is an error delivering emails, return here so we don't mark as completed
		log.Error().Err(err).Msg("could not deliver certificates to technical contact")
		return
	}

	if err = c.db.UpdateVASP(vasp); err != nil {
		log.Error().Err(err).Msg("could not update VASP email logs")
		return
	}

	// Mark certificate request as complete.
	if err = models.UpdateCertificateRequestStatus(r, models.CertificateRequestState_COMPLETED, "certificate request complete", "automated"); err != nil {
		log.Error().Err(err).Msg("could not update certificate request status")
		return
	}
	r.Status = models.CertificateRequestState_COMPLETED
	if err = c.db.UpdateCertReq(r); err != nil {
		log.Error().Err(err).Msg("could not save updated cert request")
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
	if path, err = ioutil.TempDir("", "gds_certs"); err != nil {
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

	// Iterate through the VASPs in the database.
	vasps := c.db.ListVASPs()
	defer vasps.Release()
	for vasps.Next() {
		if vasps.Error() != nil {
			log.Error().Err(err).Msg("could not iterate through database")
		}
		var vasp *pb.VASP
		if vasp, err = vasps.VASP(); err != nil {
			log.Error().Err(err).Msg("could not parse vasp record from database")
			continue
		}

		// Skip the current vasp if it is not verified, if it is verified and the
		// identity certificate is nil, log an error.
		if vasp.GetVerificationStatus() != pb.VerificationState_VERIFIED {
			continue
		} else if vasp.IdentityCertificate == nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("vasp is verified but does not have an identity certificate")
			continue
		}

		// Calculate the number of days before the VASP's certificate expires.
		notAfter := vasp.IdentityCertificate.NotAfter
		if expirationDate, err = time.Parse(time.RFC3339, notAfter); err != nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("could not parse %s's cert reissuance date")
		}
		reissuanceDate := expirationDate.Add(-time.Hour * 240) // Calculate the reissuance date, 10 days before expiration
		timeBeforeExpiration := time.Until(expirationDate)

		// TODO: handle the case were the certificate has expired, we should update the certificate record to the EXPIRED state
		switch daysBeforeExpiration := timeBeforeExpiration.Hours() / 24; {
		// Seven days before expiration, send a cert reissuance reminder to VASP.
		// NOTE: the SendContactReissuanceReminder will not send emails more than
		// once to a contact.
		case daysBeforeExpiration <= 7:
			if err = c.email.SendContactReissuanceReminder(vasp, 7, reissuanceDate); err != nil {
				log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("error sending seven day reissuance reminder")
			}

		// Ten days before expiration, reissue the VASP's identity certificate, send the email with the created pkcs12
		// password and send the whisper link, as well as notifying the TRISA admin that reissuance has started.
		case daysBeforeExpiration <= 10:
			// TODO: check that the vasps certreq is in the READY_TO_SUBMIT state to avoid the double reissue
			c.reissueIdentityCertificates(vasp)
			if _, err = c.email.SendReissuanceAdminNotification(vasp, reissuanceDate); err != nil {
				log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("error sending admin reissuance notification")
				continue
			}

		// Thirty days before expiration, send the reissuance reminder to the VASP and the TRISA admin.
		case daysBeforeExpiration <= 30:
			if err = c.email.SendContactReissuanceReminder(vasp, 30, reissuanceDate); err != nil {
				log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("error sending thirty day reissuance reminder")
			}
			if NoAdminEmailSent(vasp, 30, reissuanceDate) {
				if _, err = c.email.SendExpiresAdminNotification(vasp, reissuanceDate); err != nil {
					log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("error sending admin reissuance reminder")
					continue
				}
			}
		}
		// We need to update the vasp record in the database so that the email logs are preserved.
		if err = c.db.UpdateVASP(vasp); err != nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Msg("error updating the vasp record in the database")
		}
	}
}

// Helper function for HandleCertificateReissuance that reissues identity certificates
// for the given vasp.
func (c *CertificateManager) reissueIdentityCertificates(vasp *pb.VASP) {
	var (
		err         error
		certreq     *models.CertificateRequest
		whisperLink string
	)

	// Create a new certificate request.
	if certreq, err = models.NewCertificateRequest(vasp); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error creating certificate request for vasp %s", vasp.Id))
		return
	}

	// Update the cert req to be ready for submission.
	if err = models.UpdateCertificateRequestStatus(
		certreq,
		models.CertificateRequestState_READY_TO_SUBMIT,
		"automated certificate reissuance",
		"automated",
	); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error updating certificate request for vasp %s", vasp.Id))
		return
	}

	// Generate a new PKCS12 password
	secretType := "password"
	pkcs12password := secrets.CreateToken(16)
	whisperPasswordTemplate := "Below is the PKCS12 password which you must use to decrypt your new certificates:\n\n%s\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new secret using the secret manager.
	if err = c.secret.With(certreq.Id).CreateSecret(ctx, secretType); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error creating password secret for vasp %s", vasp.Id))
		return
	}
	if err = c.secret.With(certreq.Id).AddSecretVersion(ctx, secretType, []byte(pkcs12password)); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error creating password version for vasp %s", vasp.Id))
		return
	}

	// Using the whisper utility, create a whisper link to be sent with the ReissuanceStarted email.
	if whisperLink, err = whisper.CreateSecretLink(fmt.Sprintf(whisperPasswordTemplate, pkcs12password), "", 3, time.Now().AddDate(0, 0, 7)); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error creating whisper link for vasp %s", vasp.Id))
		return
	}

	// Send the notification email that certificate reissuance is forthcoming and provide whisper link to the PKCS12 password.
	if _, err = c.email.SendReissuanceStarted(vasp, whisperLink); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error sending reissuance started email for vasp %s", vasp.Id))
		return
	}

	// Update the certificate request in the datastore.
	if err = c.db.UpdateCertReq(certreq); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error updating certificate request for vasp %s", vasp.Id))
		return
	}

	// Save the certificate request on the VASP.
	if err = models.AppendCertReqID(vasp, certreq.Id); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error appending certificate request to vasp %s", vasp.Id))
		return
	}

	// Update the VASP information in the datastore.
	if err = c.db.UpdateVASP(vasp); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error updating vasp %s in the certman store", vasp.Id))
		return
	}
}

// Helper function for HandleCertificateReissuance checks if a certificate reissuance reminder
// has been sent to the TRISA admin within the given time window.
func NoAdminEmailSent(vasp *pb.VASP, timeWindow int, reissuanceDate time.Time) bool {
	// Make sure that the email hasn't already been sent previously.
	emailCount, err := models.GetSentAdminEmailCount(vasp, string(admin.ReissuanceReminder), timeWindow)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error retrieving admin email log for %s's reissuance reminder", vasp.Id))
		return false
	}
	if emailCount > 0 {
		return false
	}
	return true
}