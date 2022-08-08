package bff_test

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	records "github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func (s *bffTestSuite) TestCertificates() {
	require := s.Require()

	// Load fixtures for testing
	testnetCerts := &admin.ListCertificatesReply{}
	mainnetCerts := &admin.ListCertificatesReply{}
	testnetFixture := filepath.Join("testdata", "testnet", "certificates_reply.json")
	mainnetFixture := filepath.Join("testdata", "mainnet", "certificates_reply.json")
	require.NoError(loadFixture(testnetFixture, testnetCerts))
	require.NoError(loadFixture(mainnetFixture, mainnetCerts))

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err := s.client.Certificates(context.TODO())
	require.EqualError(err, "[401] this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Certificates(context.TODO())
	require.EqualError(err, "[401] user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid credentials for the remainder of the tests
	claims.Permissions = []string{"read:vasp"}
	claims.VASPs["testnet"] = authtest.TestNetVASP
	claims.VASPs["mainnet"] = authtest.MainNetVASP
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")

	// Test error message is populated when only testnet returns an error
	s.testnet.admin.UseError(mock.ListCertificatesEP, http.StatusInternalServerError, "could not retrieve testnet certificates")
	s.mainnet.admin.UseHandler(mock.ListCertificatesEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.ListCertificatesReply{
			Certificates: []admin.Certificate{},
		})
	})
	reply, err := s.client.Certificates(context.TODO())
	require.NoError(err, "expected no error when only testnet returns an error")
	require.Empty(reply.TestNet)
	require.Empty(reply.MainNet)
	require.Equal("500 Internal Server Error", reply.Error.TestNet, "expected testnet error message")
	require.Empty(reply.Error.MainNet, "expected no error when mainnet returns a valid response")

	// Test error message is populated when only mainnet returns an error
	s.testnet.admin.UseHandler(mock.ListCertificatesEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.ListCertificatesReply{
			Certificates: []admin.Certificate{},
		})
	})
	s.mainnet.admin.UseError(mock.ListCertificatesEP, http.StatusInternalServerError, "could not retrieve mainnet certificates")
	reply, err = s.client.Certificates(context.TODO())
	require.NoError(err, "expected no error when only mainnet returns an error")
	require.Empty(reply.TestNet)
	require.Empty(reply.MainNet)
	require.Equal("500 Internal Server Error", reply.Error.MainNet, "expected mainnet error message")
	require.Empty(reply.Error.TestNet, "expected no error when testnet returns a valid response")

	// Test empty results are returned even if there is no mainnet registration
	delete(claims.VASPs, "mainnet")
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")
	reply, err = s.client.Certificates(context.TODO())
	require.NoError(err, "could not retrieve certificates")
	require.Empty(reply.TestNet, "expected no testnet certificates")
	require.Empty(reply.MainNet, "expected no mainnet certificates")
	require.Empty(reply.Error, "expected no errors")

	// Test certificates are returned from both testnet and mainnet
	claims.VASPs["mainnet"] = authtest.MainNetVASP
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")
	s.testnet.admin.UseFixture(mock.ListCertificatesEP, filepath.Join("testdata", "testnet", "certificates_reply.json"))
	s.mainnet.admin.UseFixture(mock.ListCertificatesEP, filepath.Join("testdata", "mainnet", "certificates_reply.json"))
	reply, err = s.client.Certificates(context.TODO())
	require.NoError(err, "could not retrieve certificates")
	require.Len(reply.TestNet, len(testnetCerts.Certificates), "wrong number of testnet certificates")
	require.Len(reply.MainNet, len(mainnetCerts.Certificates), "wrong number of mainnet certificates")
	require.Empty(reply.Error, "expected no errors")

	// Verify the testnet certificate fields
	expected := testnetCerts.Certificates[0]
	actual := reply.TestNet[0]
	require.Equal(expected.SerialNumber, actual.SerialNumber, "expected testnet certificate serial to match")
	require.Equal(expected.IssuedAt, actual.IssuedAt, "expected testnet certificate issued date to match")
	require.Equal(expected.ExpiresAt, actual.ExpiresAt, "expected testnet certificate expiration date to match")
	require.False(actual.Revoked, "expected testnet certificate to not be revoked")
	require.Equal(expected.Details, actual.Details, "expected mainnet certificate details to match")

	// Verify one of the revoked mainnet certificates
	expected = mainnetCerts.Certificates[1]
	actual = reply.MainNet[1]
	require.Equal(expected.SerialNumber, actual.SerialNumber, "expected mainnet certificate serial to match")
	require.Equal(expected.IssuedAt, actual.IssuedAt, "expected mainnet certificate issued date to match")
	require.Equal(expected.ExpiresAt, actual.ExpiresAt, "expected mainnet certificate expiration date to match")
	require.True(actual.Revoked, "expected mainnet certificate to be revoked")
	require.Equal(expected.Details, actual.Details, "expected mainnet certificate details to match")
}

func (s *bffTestSuite) TestAttention() {
	require := s.Require()

	// Load fixtures for testing
	testnetReply := &admin.RetrieveVASPReply{}
	mainnetReply := &admin.RetrieveVASPReply{}
	testnetFixture := filepath.Join("testdata", "testnet", "retrieve_vasp_reply.json")
	mainnetFixture := filepath.Join("testdata", "mainnet", "retrieve_vasp_reply.json")
	require.NoError(loadFixture(testnetFixture, testnetReply))
	require.NoError(loadFixture(mainnetFixture, mainnetReply))

	// Create an organization in the database with no registration form
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err = s.client.Attention(context.TODO())
	require.EqualError(err, "[401] this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Attention(context.TODO())
	require.EqualError(err, "[401] user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"read:vasp"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	_, err = s.client.Attention(context.TODO())
	require.EqualError(err, "[400] missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.Attention(context.TODO())
	require.EqualError(err, "[404] no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Start registration message should be returned when there is no registration form
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	expected := &api.AttentionMessage{
		Message:  bff.StartRegistration,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_START_REGISTRATION.String(),
	}
	reply, err := s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected start registration message")
	require.Equal(expected, reply.Messages[0], "expected start registration message")

	// Start registration message should still be returned if the registration form state is empty
	org.Registration = &records.RegistrationForm{}
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected start registration message")
	require.Equal(expected, reply.Messages[0], "expected start registration message")

	// Start registration message should still be returned if the registration form has not been started
	org.Registration.State = records.NewFormState()
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected start registration message")
	require.Equal(expected, reply.Messages[0], "expected start registration message")

	// Complete registration message should be returned when the registration form has been started but not submitted
	org.Registration.State.Started = time.Now().Format(time.RFC3339)
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	expected = &api.AttentionMessage{
		Message:  bff.CompleteRegistration,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_COMPLETE_REGISTRATION.String(),
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected complete registration message")
	require.Equal(expected, reply.Messages[0], "expected complete registration message")

	// Submit mainnet message should be returned when the registration form has been submitted only to testnet
	org.Testnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	expected = &api.AttentionMessage{
		Message:  bff.SubmitMainnet,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_SUBMIT_MAINNET.String(),
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected submit mainnet message")
	require.Equal(expected, reply.Messages[0], "expected submit mainnet message")

	// Submit testnet message should be returned when the registration form has been submitted only to mainnet
	org.Testnet.Submitted = ""
	org.Mainnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	submitTestnet := &api.AttentionMessage{
		Message:  bff.SubmitTestnet,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_SUBMIT_TESTNET.String(),
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected submit testnet message")
	require.Equal(submitTestnet, reply.Messages[0], "expected submit testnet message")

	// Test an error is returned when VASP does not exist in testnet
	claims.VASPs["testnet"] = "alice0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	s.testnet.admin.UseError(mock.RetrieveVASPEP, http.StatusNotFound, "could not find VASP in database")
	_, err = s.client.Attention(context.TODO())
	require.EqualError(err, "[500] 404 Not Found", "expected error when VASP does not exist in testnet")

	// Test an error is returned when VASP does not exist in mainnet
	claims.VASPs["testnet"] = ""
	claims.VASPs["mainnet"] = "alice1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	s.mainnet.admin.UseError(mock.RetrieveVASPEP, http.StatusNotFound, "could not find VASP in database")
	_, err = s.client.Attention(context.TODO())
	require.EqualError(err, "[500] 404 Not Found", "expected error when VASP does not exist in mainnet")

	// Verify emails message should be returned when the VASP has been submitted but
	// emails are not yet verified
	vasp := &pb.VASP{}
	require.NoError(wire.Unwire(mainnetReply.VASP, vasp))
	vasp.VerificationStatus = pb.VerificationState_SUBMITTED
	data, err := wire.Rewire(vasp)
	require.NoError(err, "could not rewire VASP")
	s.mainnet.admin.UseHandler(mock.RetrieveVASPEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.RetrieveVASPReply{
			VASP: data,
		})
	})
	verifyMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.VerifyEmails, "MainNet"),
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_VERIFY_EMAILS.String(),
	}
	messages := []*api.AttentionMessage{
		submitTestnet,
		verifyMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Registration pending message should be returned when the VASP has been submitted
	// and is pending email verification
	vasp.VerificationStatus = pb.VerificationState_PENDING_REVIEW
	data, err = wire.Rewire(vasp)
	require.NoError(err, "could not rewire VASP")
	s.mainnet.admin.UseHandler(mock.RetrieveVASPEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.RetrieveVASPReply{
			VASP: data,
		})
	})
	pendingMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RegistrationPending, "MainNet"),
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_NO_ACTION.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		pendingMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Registration approved message should be returned when the VASP is verified
	vasp.VerificationStatus = pb.VerificationState_VERIFIED
	data, err = wire.Rewire(vasp)
	require.NoError(err, "could not rewire VASP")
	s.mainnet.admin.UseHandler(mock.RetrieveVASPEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.RetrieveVASPReply{
			VASP: data,
		})
	})
	approvedMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RegistrationApproved, "MainNet"),
		Severity: records.AttentionSeverity_SUCCESS.String(),
		Action:   records.AttentionAction_NO_ACTION.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		approvedMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Rejected message should be returned when the VASP state is rejected
	require.NoError(s.mainnet.admin.UseFixture(mock.RetrieveVASPEP, mainnetFixture))
	rejectMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RegistrationRejected, "MainNet"),
		Severity: records.AttentionSeverity_ALERT.String(),
		Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		rejectMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Revoked message should be returned when the certificate is revoked
	require.NoError(wire.Unwire(testnetReply.VASP, vasp))
	vasp.VerificationStatus = pb.VerificationState_VERIFIED
	vasp.IdentityCertificate.Revoked = true
	mainnetData, err := wire.Rewire(vasp)
	require.NoError(err, "could not rewire VASP")
	s.mainnet.admin.UseHandler(mock.RetrieveVASPEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.RetrieveVASPReply{
			VASP: mainnetData,
		})
	})
	revokedMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.CertificateRevoked, "MainNet"),
		Severity: records.AttentionSeverity_ALERT.String(),
		Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		revokedMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Configure testnet fixture with expired certificate
	claims.VASPs["testnet"] = "alice0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	org.Testnet.Submitted = time.Now().Format(time.RFC3339)
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	vasp = &pb.VASP{}
	require.NoError(wire.Unwire(testnetReply.VASP, vasp))
	expires := time.Now().AddDate(0, 0, 28)
	vasp.VerificationStatus = pb.VerificationState_VERIFIED
	vasp.IdentityCertificate.Revoked = false
	vasp.IdentityCertificate.NotAfter = expires.Format(time.RFC3339)
	testnetData, err := wire.Rewire(vasp)
	require.NoError(err, "could not rewire VASP")

	// Expired message should be returned when the certificate is expired
	s.testnet.admin.UseHandler(mock.RetrieveVASPEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.RetrieveVASPReply{
			VASP: testnetData,
		})
	})
	expiredTestnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RenewCertificate, "TestNet", expires.Format("January 2, 2006")),
		Severity: records.AttentionSeverity_WARNING.String(),
		Action:   records.AttentionAction_RENEW_CERTIFICATE.String(),
	}
	messages = []*api.AttentionMessage{
		expiredTestnet,
		revokedMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Should return 204 when there are no attention messages
	claims.VASPs["testnet"] = ""
	claims.VASPs["mainnet"] = ""
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Nil(reply, "expected nil reply")
}

func (s *bffTestSuite) TestRegistrationStatus() {
	require := s.Require()

	// Create an organization in the database with no directory records
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err = s.client.RegistrationStatus(context.TODO())
	require.EqualError(err, "[401] this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.RegistrationStatus(context.TODO())
	require.EqualError(err, "[401] user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"read:vasp"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	_, err = s.client.RegistrationStatus(context.TODO())
	require.EqualError(err, "[400] missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.RegistrationStatus(context.TODO())
	require.EqualError(err, "[404] no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Should return an empty response when there are no directory records
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	reply, err := s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Empty(reply, "expected empty response when there are no directory records")

	// Should return only the testnet timestamp when testnet registration has been submitted
	org.Testnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	reply, err = s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Equal(org.Testnet.Submitted, reply.TestNetSubmitted, "expected testnet timestamp to be returned")
	require.Empty(reply.MainNetSubmitted, "expected mainnet timestamp to be empty")

	// Should return only the mainnet timestamp when mainnet registration has been submitted
	org.Testnet.Submitted = ""
	org.Mainnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	reply, err = s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Equal(org.Mainnet.Submitted, reply.MainNetSubmitted, "expected mainnet timestamp to be returned")
	require.Empty(reply.TestNetSubmitted, "expected testnet timestamp to be empty")

	// Should return both timestamps when both registrations have been submitted
	org.Testnet.Submitted = time.Now().Format(time.RFC3339)
	org.Mainnet.Submitted = time.Now().Format(time.RFC3339)
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization in the database")
	reply, err = s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Equal(org.Testnet.Submitted, reply.TestNetSubmitted, "expected testnet timestamp to be returned")
	require.Equal(org.Mainnet.Submitted, reply.MainNetSubmitted, "expected mainnet timestamp to be returned")
}
