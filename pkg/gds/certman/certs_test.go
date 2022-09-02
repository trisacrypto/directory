package certman_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/certman"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/fixtures"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	trtlstore "github.com/trisacrypto/directory/pkg/gds/store/trtl"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

var (
	fixturesPath = filepath.Join("..", "testdata", "fakes.tgz")
	dbPath       = filepath.Join("testdata", "db")
)

// certTestSuite contains tests for the certificate manager to ensure that the testing
// is isolated from GDS.
type certTestSuite struct {
	suite.Suite
	fixtures *fixtures.Library
	conf     config.Config
	db       store.Store
	secret   *secrets.SecretManager
	certman  *certman.CertificateManager
}

func TestCertMan(t *testing.T) {
	suite.Run(t, new(certTestSuite))
}

func (s *certTestSuite) SetupSuite() {
	require := s.Require()

	// Discard logging to focus on test logs
	logger.Discard()

	// Initialize the fixtures library
	var err error
	s.fixtures, err = fixtures.New(fixturesPath, dbPath, fixtures.StoreTrtl)
	require.NoError(err)
}

func (s *certTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}

	s.fixtures.Close()
	logger.ResetLogger()
}

// Test that the certificate manger correctly moves certificates across the request
// pipeline.
func (s *certTestSuite) TestCertManager() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec certreq")

	// Create a secret that the certificate manager can retrieve
	sm := s.secret.With(quebecCertReq.Id)
	ctx := context.Background()
	require.NoError(sm.CreateSecret(ctx, "password"))
	require.NoError(sm.AddSecretVersion(ctx, "password", []byte("qDhAwnfMjgDEzzUC")))

	// Let the certificate manager submit the certificate request
	s.certman.HandleCertificateRequests()

	// VASP state should be changed to ISSUING_CERTIFICATE
	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)
	// Audit log should contain one additional entry for ISSUING_CERTIFICATE
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)

	// Certificate request should be updated
	certReq, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Greater(int(certReq.AuthorityId), 0)
	require.Greater(int(certReq.BatchId), 0)
	require.NotEmpty(certReq.BatchName)
	require.NotEmpty(certReq.BatchStatus)
	require.Greater(int(certReq.OrderNumber), 0)
	require.NotEmpty(certReq.CreationDate)
	require.NotEmpty(certReq.Profile)
	require.Empty(certReq.RejectReason)
	require.Equal(models.CertificateRequestState_PROCESSING, certReq.Status)
	// Audit log should contain one additional entry for PROCESSING
	require.Len(certReq.AuditLog, 3)
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, certReq.AuditLog[2].PreviousState)
	require.Equal(models.CertificateRequestState_PROCESSING, certReq.AuditLog[2].CurrentState)
	require.Equal("automated", certReq.AuditLog[2].Source)

	// Let the certificate manager process the Sectigo response
	sent := time.Now()
	s.certman.HandleCertificateRequests()

	// Secret manager should contain the certificate
	secret, err := sm.GetLatestVersion(ctx, "cert")
	require.NoError(err)
	require.NotEmpty(secret)

	// VASP should contain the new certificate
	v, err = s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	idCert := v.IdentityCertificate
	require.NotNil(idCert)
	require.Greater(int(idCert.Version), 0)
	require.NotEmpty(idCert.SerialNumber)
	require.NotEmpty(idCert.Signature)
	require.NotEmpty(idCert.SignatureAlgorithm)
	require.NotEmpty(idCert.PublicKeyAlgorithm)
	require.NotNil(idCert.Subject)
	require.NotNil(idCert.Issuer)
	_, err = time.Parse(time.RFC3339, idCert.NotBefore)
	require.NoError(err)
	_, err = time.Parse(time.RFC3339, idCert.NotAfter)
	require.NoError(err)
	require.False(idCert.Revoked)
	require.NotEmpty(idCert.Data)
	require.NotEmpty(idCert.Chain)

	// VASP should contain the certificate ID in the extra
	certIDs, err := models.GetCertIDs(v)
	require.NoError(err)
	require.Len(certIDs, 1)
	require.NotEmpty(certIDs[0])

	// VASP state should be changed to VERIFIED
	require.Equal(pb.VerificationState_VERIFIED, v.VerificationStatus)
	// Audit log should contain one additional entry for VERIFIED
	log, err = models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 6)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[5].PreviousState)
	require.Equal(pb.VerificationState_VERIFIED, log[5].CurrentState)
	require.Equal("automated", log[5].Source)

	// Certificate record should be created in the database
	cert, err := s.db.RetrieveCert(certIDs[0])
	require.NoError(err)
	require.Equal(certIDs[0], cert.Id)
	require.Equal(certReq.Id, cert.Request)
	require.Equal(v.Id, cert.Vasp)
	require.Equal(models.CertificateState_ISSUED, cert.Status)
	require.True(proto.Equal(idCert, cert.Details))

	// Email should be sent to one of the contacts
	messages := []*emails.EmailMeta{
		{
			Contact:   v.Contacts.Legal,
			To:        v.Contacts.Legal.Email,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.DeliverCertsRE,
			Reason:    "deliver_certs",
			Timestamp: sent,
		},
	}
	emails.CheckEmails(s.T(), messages)

	// Certificate request should be updated
	certReq, err = s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_COMPLETED, certReq.Status)
	require.Equal(cert.Id, certReq.Certificate)
	// Audit log should contain additional entries for DOWNLOADING, DOWNLOADED, and
	// COMPLETED
	require.Len(certReq.AuditLog, 6)
	require.Equal(models.CertificateRequestState_PROCESSING, certReq.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_DOWNLOADING, certReq.AuditLog[3].CurrentState)
	require.Equal("automated", certReq.AuditLog[3].Source)
	require.Equal(models.CertificateRequestState_DOWNLOADING, certReq.AuditLog[4].PreviousState)
	require.Equal(models.CertificateRequestState_DOWNLOADED, certReq.AuditLog[4].CurrentState)
	require.Equal("automated", certReq.AuditLog[4].Source)
	require.Equal(models.CertificateRequestState_DOWNLOADED, certReq.AuditLog[5].PreviousState)
	require.Equal(models.CertificateRequestState_COMPLETED, certReq.AuditLog[5].CurrentState)
	require.Equal("automated", certReq.AuditLog[5].Source)
}

func (s *certTestSuite) TestCertManagerThirtyDayReissuanceReminder() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Small)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()
	require := s.Require()

	// setup the datastore to contain the modified charlieVASP
	charlieVASP, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err, "could not get charlie VASP")
	charlieVASP = s.setupVASP(charlieVASP)

	// Call the certman function at 29 days, which will send
	// the thirty day cert reissuance reminder to echoVASP and
	// the TRISA admin.
	s.updateVaspIdentityCert(charlieVASP, 29)
	callTime := time.Now()
	s.certman.HandleCertificateReissuance()

	v, err := s.db.RetrieveVASP(charlieVASP.Id)
	require.NoError(err)

	// Ensure that the expected emails have been sent, using
	// the mock email client.
	messages := []*emails.EmailMeta{
		{
			To:        v.Contacts.Technical.Email,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.ReissuanceReminderRE,
			Reason:    "reissuance_reminder",
			Timestamp: callTime,
		},
		{
			To:        s.conf.Email.AdminEmail,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.ExpiresAdminNotificationRE,
			Reason:    "expires_admin_notification",
			Timestamp: callTime,
		},
	}
	emails.CheckEmails(s.T(), messages)
}

func (s *certTestSuite) TestCertManagerSevenDayReissuanceReminder() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Small)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()
	require := s.Require()

	// setup the datastore to contain the modified charlieVASP
	charlieVASP, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err, "could not get charlie VASP")
	charlieVASP = s.setupVASP(charlieVASP)

	// Call the certman function at 6 days, which will send
	// the seven day cert reissuance reminder to echoVASP.
	s.updateVaspIdentityCert(charlieVASP, 6)
	callTime := time.Now()
	s.certman.HandleCertificateReissuance()

	v, err := s.db.RetrieveVASP(charlieVASP.Id)
	require.NoError(err)

	// Ensure that the expected email has been sent, using
	// the mock email client.
	messages := []*emails.EmailMeta{
		{
			To:        v.Contacts.Technical.Email,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.ReissuanceReminderRE,
			Reason:    "reissuance_reminder",
			Timestamp: callTime,
		},
	}
	emails.CheckEmails(s.T(), messages)
}

func (s *certTestSuite) TestCertManagerReissuance() {
	require := s.Require()

	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Small)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()

	charlieVASP, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err, "could not get charlie VASP")
	charlieVASP = s.setupVASP(charlieVASP)

	// Set the second VASP in the fixtures.Small set's verification status
	// to REJECTED so that it does not get triggered for reissuance.
	deltaVASP, err := s.fixtures.GetVASP("delta")
	require.NoError(err)
	deltaVASP.VerificationStatus = pb.VerificationState_REJECTED
	require.NoError(s.db.UpdateVASP(deltaVASP))

	// Capture the number of certificate requests on the charlie VASP
	// before reissuance is triggered.
	previousReqIds, err := models.GetCertReqIDs(charlieVASP)
	require.NoError(err)
	previousNumberOfReqs := len(previousReqIds)

	// Call the certman function at 8 days, which should
	// reissue the VASP's identity certificate, send the
	// email with the created pkcs12 password and send
	// the whisper link, as well as notifying the TRISA
	// admin that reissuance has started.
	s.updateVaspIdentityCert(charlieVASP, 8)
	callTime := time.Now()
	s.certman.HandleCertificateReissuance()

	v, err := s.db.RetrieveVASP(charlieVASP.Id)
	require.NoError(err)

	reqIDs, err := models.GetCertReqIDs(v)
	require.NoError(err)
	require.Len(reqIDs, previousNumberOfReqs+1)

	// Retrieve the latest certificate request for charlie.
	certReqId := reqIDs[len(reqIDs)-1]
	certReq, err := s.db.RetrieveCertReq(certReqId)
	require.NoError(err)
	require.Equal(certReq.Status, models.CertificateRequestState_READY_TO_SUBMIT)

	// Make sure a new secret was created in the secret manager.
	sm := s.secret.With(certReq.Id)
	secret, err := sm.GetLatestVersion(context.Background(), "password")
	require.NoError(err)
	require.NotNil(secret)

	// Update the secret manager with the password that will decrypt
	// the fixture certificate used for testing, overriding the randomly
	// generated password created by updateVaspIdentityCert.
	require.NoError(sm.AddSecretVersion(context.Background(), "password", []byte("qDhAwnfMjgDEzzUC")))

	// Verify that the reissuance logic does not submit duplicate certificate requests.
	s.certman.HandleCertificateReissuance()
	v, err = s.db.RetrieveVASP(charlieVASP.Id)
	require.NoError(err)
	reqIDs, err = models.GetCertReqIDs(v)
	require.NoError(err)
	require.Len(reqIDs, previousNumberOfReqs+1, "should not have created a new certificate request")

	// Call the cert request loop once to submit the certificate request and start it's processing.
	s.certman.HandleCertificateRequests()
	v, err = s.db.RetrieveVASP(charlieVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// On the second call to the cert request loop the certificate should be downloaded and
	// attached to the VASP. The VASP should be in the VERIFIED state.
	s.certman.HandleCertificateRequests()
	v, err = s.db.RetrieveVASP(charlieVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_VERIFIED, v.VerificationStatus)

	// Ensure that the certificate request is in the COMPLETED state.
	certReq, err = s.db.RetrieveCertReq(certReqId)
	require.NoError(err)
	require.Equal(certReq.Status, models.CertificateRequestState_COMPLETED)

	// Retrieve the newly created certificate and ensure it is valid.
	idCert := v.IdentityCertificate
	require.NotNil(idCert)
	require.Greater(int(idCert.Version), 0)
	require.NotEmpty(idCert.SerialNumber)
	require.NotEmpty(idCert.Signature)
	require.NotEmpty(idCert.SignatureAlgorithm)
	require.NotEmpty(idCert.PublicKeyAlgorithm)
	require.NotNil(idCert.Subject)
	require.NotNil(idCert.Issuer)
	_, err = time.Parse(time.RFC3339, idCert.NotBefore)
	require.NoError(err)
	_, err = time.Parse(time.RFC3339, idCert.NotAfter)
	require.NoError(err)
	require.False(idCert.Revoked)
	require.NotEmpty(idCert.Data)
	require.NotEmpty(idCert.Chain)

	// Ensure that the expected email has been sent, using
	// the mock email client.
	messages := []*emails.EmailMeta{
		{
			To:        v.Contacts.Technical.Email,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.ReissuanceStartedRE,
			Reason:    "reissuance_started",
			Timestamp: callTime,
		},
		{
			To:        s.conf.Email.AdminEmail,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.ReissuanceAdminNotificationRE,
			Reason:    "reissuance_admin_notification",
			Timestamp: callTime,
		},
		{
			To:        v.Contacts.Technical.Email,
			From:      s.conf.Email.ServiceEmail,
			Subject:   emails.DeliverCertsRE,
			Reason:    "deliver_certs",
			Timestamp: callTime,
		},
	}

	// TODO: add additional testing for the email send logic in emails.getContactsToNotify()
	emails.CheckEmails(s.T(), messages)
}

func (s *certTestSuite) updateVaspIdentityCert(vasp *pb.VASP, daysUntilExpiration time.Duration) {
	days := time.Hour * 24
	daysFromNow := time.Now().Add(days * daysUntilExpiration).Format(time.RFC3339Nano)
	vasp.IdentityCertificate = &pb.Certificate{NotAfter: daysFromNow}
	s.db.UpdateVASP(vasp)
}

func (s *certTestSuite) setupVASP(vasp *pb.VASP) *pb.VASP {
	models.AddContact(vasp, "technical", &pb.Contact{
		Name:  "technical",
		Email: "technical@notmyemail.com",
	})
	models.SetContactVerification(vasp.Contacts.Technical, "", true)

	models.AddContact(vasp, "administrative", &pb.Contact{
		Name:  "administrative",
		Email: "administrative@notmyemail.com",
	})
	models.SetContactVerification(vasp.Contacts.Administrative, "", true)
	vasp.VerificationStatus = pb.VerificationState_VERIFIED

	s.db.CreateVASP(vasp)
	return vasp
}

// Test that the certificate manager rejects requests when the VASP state is invalid.
func (s *certTestSuite) TestCertManagerBadState() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()
	require := s.Require()

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec VASP")

	// Set VASP to pending review
	echoVASP.VerificationStatus = pb.VerificationState_PENDING_REVIEW
	require.NoError(s.db.UpdateVASP(echoVASP))

	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_PENDING_REVIEW, v.VerificationStatus)

	// Run the cert manager for a loop
	s.certman.HandleCertificateRequests()

	// Certificate request should be rejected before submission
	certReq, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_CR_REJECTED, certReq.Status)

	// Set VASP to rejected
	echoVASP.VerificationStatus = pb.VerificationState_REJECTED
	require.NoError(s.db.UpdateVASP(echoVASP))

	// Run the cert manager for a loop
	s.certman.HandleCertificateRequests()

	// Certificate request should be rejected before submission
	certReq, err = s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_CR_REJECTED, certReq.Status)

	// Set VASP to verified for correct submission
	echoVASP.VerificationStatus = pb.VerificationState_VERIFIED
	require.NoError(s.db.UpdateVASP(echoVASP))
	quebecCertReq.Status = models.CertificateRequestState_READY_TO_SUBMIT
	require.NoError(s.db.UpdateCertReq(quebecCertReq))

	// Move the certificate to processing
	s.certman.HandleCertificateRequests()

	// Set VASP to rejected
	echoVASP.VerificationStatus = pb.VerificationState_REJECTED
	require.NoError(s.db.UpdateVASP(echoVASP))

	// Run the cert manager for a loop
	s.certman.HandleCertificateRequests()

	// Certificate request should be rejected before download
	certReq, err = s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_CR_REJECTED, certReq.Status)
	require.Empty(certReq.Certificate)
}

// Test that the certificate manager is able to process an end entity profile.
func (s *certTestSuite) TestCertManagerEndEntityProfile() {
	s.setupCertManager(sectigo.ProfileCipherTraceEndEntityCertificate, fixtures.Full)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()
	require := s.Require()

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec VASP")

	quebecCertReq.Profile = sectigo.ProfileCipherTraceEndEntityCertificate
	quebecCertReq.Params = map[string]string{
		"organizationName":    "TRISA Member VASP",
		"localityName":        "Menlo Park",
		"stateOrProvinceName": "California",
		"countryName":         "US",
	}
	require.NoError(s.db.UpdateCertReq(quebecCertReq))

	// Create a secret that the certificate manager can retrieve.
	sm := s.secret.With(quebecCertReq.Id)
	ctx := context.Background()
	require.NoError(sm.CreateSecret(ctx, "password"))
	require.NoError(sm.AddSecretVersion(ctx, "password", []byte("qDhAwnfMjgDEzzUC")))

	// Run the certificate manager through two iterations to fully process the request.
	s.certman.HandleCertificateRequests()
	s.certman.HandleCertificateRequests()

	// VASP should contain the new certificate
	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_VERIFIED, v.VerificationStatus)
	require.NotNil(v.IdentityCertificate)

	// Certificate request should be updated
	cert, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_COMPLETED, cert.Status)
}

// Test that the certificate manager is able to process a CipherTraceEE profile.
func (s *certTestSuite) TestCertManagerCipherTraceEEProfile() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()
	require := s.Require()

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec VASP")

	quebecCertReq.Profile = sectigo.ProfileCipherTraceEE
	require.NoError(s.db.UpdateCertReq(quebecCertReq))

	// Create a secret that the certificate manager can retrieve
	sm := s.secret.With(quebecCertReq.Id)
	ctx := context.Background()
	require.NoError(sm.CreateSecret(ctx, "password"))
	require.NoError(sm.AddSecretVersion(ctx, "password", []byte("qDhAwnfMjgDEzzUC")))

	// Run the certificate manager through two iterations to fully process the request.
	s.certman.HandleCertificateRequests()
	s.certman.HandleCertificateRequests()

	// VASP should contain the new certificate
	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_VERIFIED, v.VerificationStatus)
	require.NotNil(v.IdentityCertificate)

	// Certificate request should be updated
	cert, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_COMPLETED, cert.Status)
}

// Test that certificate submission fails if the user available balance is 0.
func (s *certTestSuite) TestSubmitNoBalance() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	mock.Handle(sectigo.AuthorityUserBalanceAvailableEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, 0)
	})

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec VASP")

	// Run the CertManager for a tick
	s.certman.HandleCertificateRequests()

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Cert request should still be in the READY_TO_SUBMIT state
	cert, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.Status)

	// Audit log should be updated
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)
}

// Test that the certificate submission fails if there is no available password.
func (s *certTestSuite) TestSubmitNoPassword() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec VASP")

	// Run the CertManager for a tick
	s.certman.HandleCertificateRequests()

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Cert request should still be in the READY_TO_SUBMIT state
	cert, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.Status)

	// Audit log should be updated
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)
}

// Test that the certificate submission fails if the batch request fails.
func (s *certTestSuite) TestSubmitBatchError() {
	s.setupCertManager(sectigo.ProfileCipherTraceEndEntityCertificate, fixtures.Full)
	defer s.teardownCertManager()
	defer s.fixtures.LoadReferenceFixtures()
	require := s.Require()

	echoVASP, err := s.fixtures.GetVASP("echo")
	require.NoError(err, "could not get echo VASP")
	quebecCertReq, err := s.fixtures.GetCertReq("quebec")
	require.NoError(err, "could not get quebec VASP")

	// Create a secret that the certificate manager can retrieve
	sm := s.secret.With(quebecCertReq.Id)
	ctx := context.Background()
	require.NoError(sm.CreateSecret(ctx, "password"))
	require.NoError(sm.AddSecretVersion(ctx, "password", []byte("qDhAwnfMjgDEzzUC")))

<<<<<<< HEAD
	// Certificate request with a missing country name
=======
	// Create a valid certificate request with extended parameters
>>>>>>> origin/main
	quebecCertReq.Params = map[string]string{
		"organizationName":    "TRISA Member VASP",
		"localityName":        "Menlo Park",
		"stateOrProvinceName": "California",
<<<<<<< HEAD
	}
	require.NoError(s.db.UpdateCertReq(quebecCertReq))

=======
		"country":             "US",
	}
	require.NoError(s.db.UpdateCertReq(quebecCertReq))

	// Ensure that Sectigo returns an error response when the batch is submitted.
	mock.Handle(sectigo.CreateSingleCertBatchEP, func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

>>>>>>> origin/main
	// Run the CertManager for a tick
	s.certman.HandleCertificateRequests()

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.db.RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Cert request should still be in the READY_TO_SUBMIT state
	cert, err := s.db.RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
<<<<<<< HEAD
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.Status)
=======
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.Status, "certificate request is not in ready to submit state")
>>>>>>> origin/main

	// Audit log should be updated
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)
}

// Test that the certificate processing fails if the batch status request fails.
func (s *certTestSuite) TestProcessBatchDetailError() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	foxtrot, err := s.fixtures.GetVASP("foxtrot")
	require.NoError(err, "could not get foxtrot VASP")

	// Batch detail returns an error
	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	// Run cert manager for one loop
	s.certman.HandleCertificateRequests()
	require.NoError(err, "certman loop unsuccessful")

	v, err := s.db.RetrieveVASP(foxtrot.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Batch status can't be retrieved from both the detail and status endpoints.
	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.BatchResponse{
			BatchID:      42,
			CreationDate: time.Now().Format(time.RFC3339),
		})
	})
	mock.Handle(sectigo.BatchStatusEP, func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	// Run cert manager for one loop
	s.certman.HandleCertificateRequests()

	v, err = s.db.RetrieveVASP(foxtrot.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)
}

// Test that the certificate processing fails if there is still an active batch.
func (s *certTestSuite) TestProcessActiveBatch() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	foxtrot, err := s.fixtures.GetVASP("foxtrot")
	require.NoError(err, "could not get foxtrot VASP")
	sierra, err := s.fixtures.GetCertReq("sierra")
	require.NoError(err, "could not get sierra VASP")

	// Batch detail returns an error
	mock.Handle(sectigo.BatchProcessingInfoEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.ProcessingInfoResponse{
			Active:  1,
			Success: 0,
			Failed:  0,
		})
	})

	// Run cert manager for one loop
	s.certman.HandleCertificateRequests()

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.db.RetrieveVASP(foxtrot.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to PROCESSING
	cert, err := s.db.RetrieveCertReq(sierra.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[2].Source)
}

// Test that the certificate processing fails if the batch request is rejected.
func (s *certTestSuite) TestProcessRejected() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	foxtrot, err := s.fixtures.GetVASP("foxtrot")
	require.NoError(err, "could not get foxtrot VASP")
	sierra, err := s.fixtures.GetCertReq("sierra")
	require.NoError(err, "could not get sierra VASP")

	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.BatchResponse{
			BatchID:      42,
			CreationDate: time.Now().Format(time.RFC3339),
			Status:       sectigo.BatchStatusRejected,
		})
	})
	mock.Handle(sectigo.BatchProcessingInfoEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.ProcessingInfoResponse{
			Active:  0,
			Success: 0,
			Failed:  1,
		})
	})

	// Run cert manager for one loop
	s.certman.HandleCertificateRequests()

	// VASP state should be still be ISSUING_CERTIFICATE
	v, err := s.db.RetrieveVASP(foxtrot.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to CR_REJECTED
	cert, err := s.db.RetrieveCertReq(sierra.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_CR_REJECTED, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_CR_REJECTED, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[3].Source)
}

// Test that the certificate processing fails if the batch request errors.
func (s *certTestSuite) TestProcessBatchError() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	foxtrot, err := s.fixtures.GetVASP("foxtrot")
	require.NoError(err, "could not get foxtrot VASP")
	sierra, err := s.fixtures.GetCertReq("sierra")
	require.NoError(err, "could not get sierra VASP")

	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.BatchResponse{
			BatchID:      42,
			CreationDate: time.Now().Format(time.RFC3339),
			Status:       sectigo.BatchStatusNotAcceptable,
		})
	})
	mock.Handle(sectigo.BatchProcessingInfoEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.ProcessingInfoResponse{
			Active:  0,
			Success: 0,
			Failed:  1,
		})
	})

	// Run cert manager for one loop
	s.certman.HandleCertificateRequests()

	// VASP state should be still be ISSUING_CERTIFICATE
	v, err := s.db.RetrieveVASP(foxtrot.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to CR_ERRORED
	cert, err := s.db.RetrieveCertReq(sierra.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_CR_ERRORED, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_CR_ERRORED, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[3].Source)
}

// Test that the certificate processing fails if the batch processing info request
// returns an unhandled sectigo state.
func (s *certTestSuite) TestProcessBatchNoSuccess() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	require := s.Require()

	foxtrot, err := s.fixtures.GetVASP("foxtrot")
	require.NoError(err, "could not get foxtrot VASP")
	sierra, err := s.fixtures.GetCertReq("sierra")
	require.NoError(err, "could not get sierra VASP")

	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.BatchResponse{
			BatchID:      42,
			CreationDate: time.Now().Format(time.RFC3339),
			Status:       sectigo.BatchStatusNotAcceptable,
		})
	})

	// Run cert manager for one loop
	s.certman.HandleCertificateRequests()

	// VASP state should be still be ISSUING_CERTIFICATE
	v, err := s.db.RetrieveVASP(foxtrot.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to PROCESSING
	cert, err := s.db.RetrieveCertReq(sierra.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[3].Source)
}

func (s *certTestSuite) TestCertManagerRequestLoop() {
	s.setupCertManager(sectigo.ProfileCipherTraceEE, fixtures.Full)
	defer s.teardownCertManager()
	s.runCertManager(s.conf.CertMan.RequestInterval)
}

func (s *certTestSuite) setupCertManager(profile string, fType fixtures.FixtureType) {
	require := s.Require()

	// Load fixtures into the library
	require.NoError(s.fixtures.Load(fType))

	// Get mock configuration values
	s.conf = gds.MockConfig()

	// Create the certificate manager configuration
	var err error
<<<<<<< HEAD
	certPath, err := ioutil.TempDir("testdata", "certs-*")
=======
	certPath, err := os.MkdirTemp("testdata", "certs-*")
>>>>>>> origin/main
	require.NoError(err, "could not create cert storage")
	s.conf.CertMan.Storage = certPath
	s.conf.CertMan.RequestInterval = time.Millisecond
	s.conf.CertMan.Sectigo.Profile = profile

	// Initialize the trtl store
	conn, err := s.fixtures.ConnectTrtl(context.Background())
	require.NoError(err, "could not connect to trtl database")
	s.db, err = trtlstore.NewMock(conn)
	require.NoError(err, "could not create trtl store")

	// Initialize the secret manager
	s.secret, err = secrets.NewMock(s.conf.Secrets)
	require.NoError(err, "could not create secret manager")

	// Initialize the email manager
	email, err := emails.New(s.conf.Email)
	require.NoError(err, "could not create email manager")

	// Initialize the certificate manager
	require.NoError(os.MkdirAll(s.conf.CertMan.Storage, 0755))
	s.certman, err = certman.New(s.conf.CertMan, s.db, s.secret, email)
	require.NoError(err, "could not create certificate manager")
}

func (s *certTestSuite) teardownCertManager() {
	require := s.Require()
	emails.PurgeMockEmails()
	s.db.Close()
	s.fixtures.Reset()
	require.NoError(os.RemoveAll(s.conf.CertMan.Storage))
}

// Helper function that spins up the CertificateManager for the specified duration,
// sends the stop signal, and waits for it to finish.
func (s *certTestSuite) runCertManager(requestInterval time.Duration) {
	// Start the certificate manager
	wg := sync.WaitGroup{}
	s.certman.Run(&wg)

	// Wait for the interval to elapse
	time.Sleep(requestInterval)

	// Make sure that the certificate manager is stopped before we proceed
	s.certman.Stop()
	wg.Wait()
}
