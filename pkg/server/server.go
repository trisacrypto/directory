package server

import "github.com/trisacrypto/directory/pkg/config"

type Server struct {
	conf *config.Config
}

func New() (*Server, error) {
	conf, err := config.New()
	if err != nil {
		return nil, err
	}
	return &Server{conf: conf}, nil
}

func (s *Server) Serve() error {
	return nil
}
