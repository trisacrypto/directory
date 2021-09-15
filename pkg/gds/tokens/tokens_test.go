package tokens_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

type TokenTestSuite struct {
	suite.Suite
	testdata map[string]string
}

func (suite *TokenTestSuite) SetupSuite() {
	// Create the keys map from the testdata directory to create new token managers.
	suite.testdata = make(map[string]string)
	suite.testdata["1yAwhf28bXi3IWP6FYcGa0dcrfq"] = "testdata/1yAwhf28bXi3IWP6FYcGa0dcrfq.pem"
	suite.testdata["1yAxs5vPqCrg433fPrFENevvzen"] = "testdata/1yAxs5vPqCrg433fPrFENevvzen.pem"
}

func (suite *TokenTestSuite) TestTokenManager() {
	tm, err := tokens.New(suite.testdata)
	suite.NoError(err, "could not initialize token manager")

	keys := tm.Keys()
	suite.Len(keys, 2)
	suite.Equal("1yAxs5vPqCrg433fPrFENevvzen", tm.CurrentKey().String())
}

// Execute suite as a go test.
func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
