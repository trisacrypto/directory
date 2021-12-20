package gds_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Test that the certificate manger correctly moves certificates across the request
// pipeline.
func (s *gdsTestSuite) TestCertManager() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	echoVASP := s.fixtures[vasps]["echo"].(*pb.VASP)
	quebecCertReq := s.fixtures[certreqs]["quebec"].(*models.CertificateRequest)

	// Create a secret that the certificate manager can retrieve
	sm := s.svc.GetSecretManager().With(quebecCertReq.Id)
	ctx := context.Background()
	require.NoError(sm.CreateSecret(ctx, "password"))
	require.NoError(sm.AddSecretVersion(ctx, "password", []byte("qDhAwnfMjgDEzzUC")))

	// Let the certificate manager submit the certificate request
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP state should be changed to ISSUING_CERTIFICATE
	v, err := s.svc.GetStore().RetrieveVASP(echoVASP.Id)
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
	cert, err := s.svc.GetStore().RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Greater(int(cert.AuthorityId), 0)
	require.Greater(int(cert.BatchId), 0)
	require.NotEmpty(cert.BatchName)
	require.NotEmpty(cert.BatchStatus)
	require.Greater(int(cert.OrderNumber), 0)
	require.NotEmpty(cert.CreationDate)
	require.NotEmpty(cert.Profile)
	require.Empty(cert.RejectReason)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.Status)
	// Audit log should contain one additional entry for PROCESSING
	require.Len(cert.AuditLog, 3)
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.AuditLog[2].PreviousState)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[2].CurrentState)
	require.Equal("automated", cert.AuditLog[2].Source)

	// Let the certificate manager process the Sectigo response
	sent := time.Now()
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// Wait for the download routine to finish
	time.Sleep(time.Second)

	// Secret manager should contain the certificate
	secret, err := sm.GetLatestVersion(ctx, "cert")
	require.NoError(err)
	require.NotEmpty(secret)

	// VASP should contain the new certificate
	v, err = s.svc.GetStore().RetrieveVASP(echoVASP.Id)
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

	// VASP state should be changed to VERIFIED
	require.Equal(pb.VerificationState_VERIFIED, v.VerificationStatus)
	// Audit log should contain one additional entry for VERIFIED
	log, err = models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 6)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[5].PreviousState)
	require.Equal(pb.VerificationState_VERIFIED, log[5].CurrentState)
	require.Equal("automated", log[5].Source)

	// Email should be sent to one of the contacts
	messages := []*emailMeta{
		{
			contact:   v.Contacts.Legal,
			to:        v.Contacts.Legal.Email,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.DeliverCertsRE,
			reason:    "deliver_certs",
			timestamp: sent,
		},
	}
	s.CheckEmails(messages)

	// Certificate request should be updated
	cert, err = s.svc.GetStore().RetrieveCertReq(quebecCertReq.Id)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_COMPLETED, cert.Status)
	// Audit log should contain additional entries for DOWNLOADING, DOWNLOADED, and
	// COMPLETED
	require.Len(cert.AuditLog, 6)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_DOWNLOADING, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[3].Source)
	require.Equal(models.CertificateRequestState_DOWNLOADING, cert.AuditLog[4].PreviousState)
	require.Equal(models.CertificateRequestState_DOWNLOADED, cert.AuditLog[4].CurrentState)
	require.Equal("automated", cert.AuditLog[4].Source)
	require.Equal(models.CertificateRequestState_DOWNLOADED, cert.AuditLog[5].PreviousState)
	require.Equal(models.CertificateRequestState_COMPLETED, cert.AuditLog[5].CurrentState)
	require.Equal("automated", cert.AuditLog[5].Source)
}

// Test that certificate submission fails if the user available balance is 0.
func (s *gdsTestSuite) TestSubmitNoBalance() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	mock.Handle(sectigo.AuthorityUserBalanceAvailableEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, 0)
	})

	echoVASP := s.fixtures[vasps]["echo"].(*pb.VASP)

	// Run the CertManager for a tick
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.svc.GetStore().RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Audit log should be updated
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)
}

// Test that the certificate submission fails if there is no available password.
func (s *gdsTestSuite) TestSubmitNoPassword() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	echoVASP := s.fixtures[vasps]["echo"].(*pb.VASP)

	// Run the CertManager for a tick
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.svc.GetStore().RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Audit log should be updated
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)
}

// Test that the certificate submission fails if the batch request fails.
func (s *gdsTestSuite) TestSubmitBatchError() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	echoVASP := s.fixtures[vasps]["echo"].(*pb.VASP)
	quebecCertReq := s.fixtures[certreqs]["quebec"].(*models.CertificateRequest)

	// Create a secret that the certificate manager can retrieve
	sm := s.svc.GetSecretManager().With(quebecCertReq.Id)
	ctx := context.Background()
	require.NoError(sm.CreateSecret(ctx, "password"))
	require.NoError(sm.AddSecretVersion(ctx, "password", []byte("qDhAwnfMjgDEzzUC")))

	mock.Handle(sectigo.CreateSingleCertBatchEP, func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	// Run the CertManager for a tick
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.svc.GetStore().RetrieveVASP(echoVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Audit log should be updated
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 5)
	require.Equal(pb.VerificationState_REVIEWED, log[4].PreviousState)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, log[4].CurrentState)
	require.Equal("automated", log[4].Source)
}

// Test that the certificate processing fails if the batch status request fails.
func (s *gdsTestSuite) TestProcessBatchDetailError() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	foxtrotId := s.fixtures[vasps]["foxtrot"].(*pb.VASP).Id

	// Batch detail returns an error
	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	s.runCertManager(s.svc.GetConf().CertMan.Interval)
	v, err := s.svc.GetStore().RetrieveVASP(foxtrotId)
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
	s.runCertManager(s.svc.GetConf().CertMan.Interval)
	v, err = s.svc.GetStore().RetrieveVASP(foxtrotId)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)
}

// Test that the certificate processing fails if there is still an active batch.
func (s *gdsTestSuite) TestProcessActiveBatch() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	foxtrotId := s.fixtures[vasps]["foxtrot"].(*pb.VASP).Id
	sierraId := s.fixtures[certreqs]["sierra"].(*models.CertificateRequest).Id

	// Batch detail returns an error
	mock.Handle(sectigo.BatchProcessingInfoEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.ProcessingInfoResponse{
			Active:  1,
			Success: 0,
			Failed:  0,
		})
	})
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP should still be in the ISSUING_CERTIFICATE state
	v, err := s.svc.GetStore().RetrieveVASP(foxtrotId)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to PROCESSING
	cert, err := s.svc.GetStore().RetrieveCertReq(sierraId)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[2].Source)
}

// Test that the certificate processing fails if the batch request is rejected.
func (s *gdsTestSuite) TestProcessRejected() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	foxtrotId := s.fixtures[vasps]["foxtrot"].(*pb.VASP).Id
	sierraId := s.fixtures[certreqs]["sierra"].(*models.CertificateRequest).Id

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
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP state should be still be ISSUING_CERTIFICATE
	v, err := s.svc.GetStore().RetrieveVASP(foxtrotId)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to CR_REJECTED
	cert, err := s.svc.GetStore().RetrieveCertReq(sierraId)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_CR_REJECTED, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_CR_REJECTED, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[3].Source)
}

// Test that the certificate processing fails if the batch request errors.
func (s *gdsTestSuite) TestProcessBatchError() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	foxtrotId := s.fixtures[vasps]["foxtrot"].(*pb.VASP).Id
	sierraId := s.fixtures[certreqs]["sierra"].(*models.CertificateRequest).Id

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
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP state should be still be ISSUING_CERTIFICATE
	v, err := s.svc.GetStore().RetrieveVASP(foxtrotId)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to CR_ERRORED
	cert, err := s.svc.GetStore().RetrieveCertReq(sierraId)
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
func (s *gdsTestSuite) TestProcessBatchNoSuccess() {
	s.setupCertManager()
	defer s.teardownCertManager()
	require := s.Require()

	foxtrotId := s.fixtures[vasps]["foxtrot"].(*pb.VASP).Id
	sierraId := s.fixtures[certreqs]["sierra"].(*models.CertificateRequest).Id

	mock.Handle(sectigo.BatchDetailEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &sectigo.BatchResponse{
			BatchID:      42,
			CreationDate: time.Now().Format(time.RFC3339),
			Status:       sectigo.BatchStatusNotAcceptable,
		})
	})
	s.runCertManager(s.svc.GetConf().CertMan.Interval)

	// VASP state should be still be ISSUING_CERTIFICATE
	v, err := s.svc.GetStore().RetrieveVASP(foxtrotId)
	require.NoError(err)
	require.Equal(pb.VerificationState_ISSUING_CERTIFICATE, v.VerificationStatus)

	// Certificate request state should be changed to PROCESSING
	cert, err := s.svc.GetStore().RetrieveCertReq(sierraId)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.Status)

	// Audit log should be updated
	require.Len(cert.AuditLog, 4)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].PreviousState)
	require.Equal(models.CertificateRequestState_PROCESSING, cert.AuditLog[3].CurrentState)
	require.Equal("automated", cert.AuditLog[3].Source)
}

func (s *gdsTestSuite) setupCertManager() {
	require := s.Require()
	tmp, err := ioutil.TempDir("testdata", "certs-*")
	require.NoError(err)
	conf := gds.MockConfig()
	conf.CertMan = config.CertManConfig{
		Interval: time.Millisecond,
		Storage:  tmp,
	}
	require.NoError(os.MkdirAll(conf.CertMan.Storage, 0755))
	s.SetConfig(conf)
	s.LoadFullFixtures()
}

func (s *gdsTestSuite) teardownCertManager() {
	s.ResetConfig()
	s.ResetFullFixtures()
	emails.PurgeMockEmails()
	os.RemoveAll(s.svc.GetConf().CertMan.Storage)
}

// Helper function that spins up the CertificateManager for the specified duration,
// sends the stop signal, and waits for it to finish.
func (s *gdsTestSuite) runCertManager(interval time.Duration) {
	// Start the certificate manager
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.svc.CertManager(stop)
	}()

	// Wait for the interval to elapse
	time.Sleep(interval)

	// Make sure that the certificate manager is stopped before we proceed
	close(stop)
	wg.Wait()
}
