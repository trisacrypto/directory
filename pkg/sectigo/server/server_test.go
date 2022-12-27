package server_test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/server"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

type serverTestSuite struct {
	suite.Suite
	conf    server.Config
	profile sectigo.Config
	srv     *server.Server
	client  *sectigo.Sectigo
}

func (s *serverTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	conf := server.Config{
		BindAddr:   "127.0.0.1:0",
		Mode:       gin.TestMode,
		LogLevel:   logger.LevelDecoder(zerolog.PanicLevel),
		ConsoleLog: false,
		Auth: server.AuthConfig{
			Username: sectigo.MockUsername,
			Password: sectigo.MockPassword,
			Issuer:   "http://localhost",
			Subject:  "testuser",
			Scopes:   []string{"ROLE_USER"},
		},
	}
	s.conf, err = conf.Mark()
	require.NoError(err, "could not validate server configuration")

	s.srv, err = server.New(s.conf)
	require.NoError(err, "could not create sectigo server")
	go s.srv.Serve()
	time.Sleep(500 * time.Millisecond)

	s.profile = sectigo.Config{
		Username: sectigo.MockUsername,
		Password: sectigo.MockPassword,
		Profile:  sectigo.ProfileCipherTraceEE,
		Testing:  true,
		Endpoint: s.srv.URL(),
	}
	s.client, err = sectigo.New(s.profile)
	require.NoError(err, "could not create sectigo client")

	// Ensure the client is authenticated
	err = s.client.Authenticate()
	require.NoError(err, "could not authenticate client")
}

func (s *serverTestSuite) TearDownSuite() {
	require := s.Require()
	err := s.srv.Shutdown()
	require.NoError(err)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}
