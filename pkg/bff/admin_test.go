package bff_test

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
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

	// Test an error is returned when only testnet returns an error
	s.testnet.admin.UseError(mock.ListCertificatesEP, http.StatusInternalServerError, "could not retrieve testnet certificates")
	s.mainnet.admin.UseHandler(mock.ListCertificatesEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.ListCertificatesReply{
			Certificates: []admin.Certificate{},
		})
	})
	reply, err := s.client.Certificates(context.TODO())
	require.Empty(reply.TestNet)
	require.Empty(reply.MainNet)
	require.Equal("500 Internal Server Error", reply.Error.TestNet, "expected error when testnet returns an error")
	require.Empty(reply.Error.MainNet, "expected no error when mainnet returns a valid response")

	// Test an error is returned when only mainnet returns an error
	s.testnet.admin.UseHandler(mock.ListCertificatesEP, func(c *gin.Context) {
		c.JSON(http.StatusOK, &admin.ListCertificatesReply{
			Certificates: []admin.Certificate{},
		})
	})
	s.mainnet.admin.UseError(mock.ListCertificatesEP, http.StatusInternalServerError, "could not retrieve mainnet certificates")
	reply, err = s.client.Certificates(context.TODO())
	require.Empty(reply.TestNet)
	require.Empty(reply.MainNet)
	require.Equal("500 Internal Server Error", reply.Error.MainNet, "expected error when mainnet returns an error")
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
