/*
Package clive provides CLI-Live interactions with Auth0 by running a local server for
OAuth challenges and handling them on behalf of the user.
*/
package clive

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const BindAddr = "127.0.0.1:4784"

type Server struct {
	conf  Config
	srv   *http.Server
	echan chan error
}

func New(conf Config) (*Server, error) {
	s := &Server{conf: conf, echan: make(chan error, 1)}

	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.png", s.Favicon)
	mux.HandleFunc("/auth/callback", s.Authorize)

	s.srv = &http.Server{
		Addr:         BindAddr,
		Handler:      mux,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s, nil
}

func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	var sock net.Listener
	if sock, err = net.Listen("tcp", BindAddr); err != nil {
		return err
	}

	go func() {
		if err := s.srv.Serve(sock); err != nil && err != http.ErrServerClosed {
			s.echan <- err
		}
	}()

	// Listen for any errors that might have occurred and wait for all go routines to stop
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() (err error) {
	// Require shutdown in 5 seconds without blocking
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) Authorize(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(MustAsset("index.html"))
}

func (s *Server) Favicon(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(MustAsset("favicon.png"))
}
