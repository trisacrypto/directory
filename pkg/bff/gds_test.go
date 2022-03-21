package bff_test

import (
	"context"
	"fmt"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
)

func (s *bffTestSuite) TestLookup() {
	require := s.Require()
	params := &api.LookupParams{}

	// Test Bad Request (no parameters)
	_, err := s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[400] must provide either uuid or common_name in query params", "expected a 400 error with no params")

	// Provide some params
	params.CommonName = "api.alice.vaspbot.net"

	// Test NotFound returns a 404
	require.NoError(s.testnet.UseError(mock.LookupRPC, codes.NotFound, "tesnet not found"))
	require.NoError(s.mainnet.UseError(mock.LookupRPC, codes.NotFound, "mainnet not found"))
	_, err = s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[404] no results returned for query", "expected a 404 error when both GDSes return not found")

	// Test InternalError when both GDSes return Unavailable
	require.NoError(s.testnet.UseError(mock.LookupRPC, codes.Unavailable, "testnet cannot connect"))
	require.NoError(s.mainnet.UseError(mock.LookupRPC, codes.Unavailable, "mainnet cannot connect"))
	_, err = s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[500] unable to execute Lookup request", "expected a 500 error when both GDSes return unavailable")

	// Test one result from TestNet
	require.NoError(s.testnet.UseFixture(mock.LookupRPC, "testdata/testnet/lookup_reply.json"))
	require.NoError(s.mainnet.UseError(mock.LookupRPC, codes.NotFound, "mainnet not found"))
	rep, err := s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from testnet")
	require.Len(rep.Results, 1, "expected one result back from server")
	require.Equal("6a57fea4-8fb7-42f3-bf0c-55fecccd2e53", rep.Results[0]["id"])

	// Test one result from MainNet
	require.NoError(s.testnet.UseError(mock.LookupRPC, codes.NotFound, "testnet not found"))
	require.NoError(s.mainnet.UseFixture(mock.LookupRPC, "testdata/mainnet/lookup_reply.json"))
	rep, err = s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from mainnet")
	require.Len(rep.Results, 1, "expected one result back from server")
	require.Equal("ca0cff66-719f-4a62-8086-be953699b27d", rep.Results[0]["id"])

	// Test results from both TestNet and MainNet
	require.NoError(s.testnet.UseFixture(mock.LookupRPC, "testdata/testnet/lookup_reply.json"))
	require.NoError(s.mainnet.UseFixture(mock.LookupRPC, "testdata/mainnet/lookup_reply.json"))
	rep, err = s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from mainnet and testnet")
	require.Len(rep.Results, 2, "expected two results back from server")
}

func (s *bffTestSuite) TestRegister() {
	require := s.Require()

	// Test both the testnet and the mainnet registration
	for _, network := range []string{"testnet", "mainnet"} {
		req := &api.RegisterRequest{
			Network: network,
		}

		// Test Errors first - should make no calls to the mock GDS because the input is invalid.
		// Test Business Category is required
		_, err := s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] business category is required")

		// Test Entity is required
		req.BusinessCategory = "FOO"
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] entity is required")

		// Test Contacts are required
		req.Entity = map[string]interface{}{"name": 1}
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] contacts are required")

		// Test TRIXO is required
		req.Contacts = map[string]interface{}{"technical": "red"}
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] trixo is required")

		// Test entity must be valid
		req.TRIXO = map[string]interface{}{"primary_regulator": 1}
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] could not parse legal person entity")

		// Test contacts must be valid
		req.Entity, err = loadFixture("testdata/entity.json")
		require.NoError(err, "could not load testdata/entity.json")
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] could not parse contacts")

		// Test business category must be valid
		req.Contacts, err = loadFixture("testdata/contacts.json")
		require.NoError(err, "could not load testdata/contacts.json")
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] could not parse \"FOO\" into a business category")

		// Test TRIXO must be valid
		req.BusinessCategory = models.BusinessCategoryPrivate.String()
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[400] could not parse TRIXO form")

		// We now have a valid request from BFF's perspective because BFF only handles
		// the intermediate parsing of the protocol buffers. However, the GDS can still
		// validate e.g. if the common name doesn't match the endpoint or the entity is
		// not a valid IVMS 101 LegalPerson struct. So we simulate the GDS returning an
		// invalid argument, which the BFF should pass back as a 400 error.
		req.TRIXO, err = loadFixture("testdata/trixo.json")
		require.NoError(err, "could not load testdata/trixo.json")
		req.TRISAEndpoint = "trisa.example.com:443"
		req.CommonName = "trisa.example.com"
		req.Website = "https://example.com"
		req.VASPCategories = []string{models.VASPCategoryKiosk, models.VASPCategoryProject}
		req.EstablishedOn = "2019-06-21"

		// Identify the mock being used in this loop
		var mgds *mock.GDS
		switch network {
		case "testnet":
			mgds = s.testnet
		case "mainnet":
			mgds = s.mainnet
		}

		// Reset the calls on the mocks to ensure the correct mock GDS is being called
		expectedCalls := make(map[string]int)
		s.testnet.Reset()
		s.mainnet.Reset()

		// Test Invalid Argument Error
		mgds.UseError(mock.RegisterRPC, codes.InvalidArgument, "the TRISA endpoint is not valid")
		_, err = s.client.Register(context.TODO(), req)
		expectedCalls[network]++
		require.EqualError(err, "[400] the TRISA endpoint is not valid")
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test Already Exists error
		mgds.UseError(mock.RegisterRPC, codes.AlreadyExists, "this VASP is already registered")
		_, err = s.client.Register(context.TODO(), req)
		expectedCalls[network]++
		require.EqualError(err, "[400] this VASP is already registered")
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test Aborted error
		mgds.UseError(mock.RegisterRPC, codes.Aborted, "a conflict occurred")
		_, err = s.client.Register(context.TODO(), req)
		expectedCalls[network]++
		require.EqualError(err, "[409] a conflict occurred")
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test Timeout error
		mgds.UseError(mock.RegisterRPC, codes.DeadlineExceeded, "deadline exceeded")
		_, err = s.client.Register(context.TODO(), req)
		expectedCalls[network]++
		require.EqualError(err, fmt.Sprintf("[500] could not register with %s", network))
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test FailedPrecondition error
		mgds.UseError(mock.RegisterRPC, codes.FailedPrecondition, "couldn't access database")
		_, err = s.client.Register(context.TODO(), req)
		expectedCalls[network]++
		require.EqualError(err, fmt.Sprintf("[500] could not register with %s", network))
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test a valid register reply
		err = mgds.UseFixture(mock.RegisterRPC, fmt.Sprintf("testdata/%s/register_reply.json", network))
		require.NoError(err, "could not load register reply fixture")

		rep, err := s.client.Register(context.TODO(), req)
		expectedCalls[network]++
		require.NoError(err, "could not make register call with valid payload")
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Check the register response is valid
		require.Empty(rep.Error, "an error message was returned from the server")
		require.NotEmpty(rep.Id, "the ID was not returned from the server")
		require.NotEmpty(rep.RegisteredDirectory, "the registered directory was not returned from the server")
		require.Equal(rep.CommonName, "trisa.example.com", "the common name was not returned from the server")
		require.Equal(rep.Status, "PENDING_REVIEW", "the verification status was not returned by the server")
		require.Equal(rep.Message, "thank you for registering", "a message was not returned from the server")
		require.Equal(rep.PKCS12Password, "supersecret", "a pkcs12 password was not returned from the server")

		// Test that a post to an incorrect network returns an error.
		req.Network = "foo"
		_, err = s.client.Register(context.TODO(), req)
		require.EqualError(err, "[404] network should be either testnet or mainnet")
		require.Equal(expectedCalls["testnet"], s.testnet.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)
	}
}
