/*
Package clive provides CLI-Live interactions with Auth0 by running a local server for
OAuth challenges and handling them on behalf of the user.
*/
package clive

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	BindAddr    = "127.0.0.1:4784"
	RedirectURI = "http://localhost:4784/auth/callback"
)

type Server struct {
	conf         Config
	srv          *http.Server
	echan        chan error
	codeVerifier string
	state        string
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
	ctx := RenderContext{}

	// Parse query params from request
	params := req.URL.Query()
	if params.Has("error") {
		// Begin error handling
		ctx.Error = params.Get("error")
		ctx.Description = params.Get("error_description")
	}

	if params.Has("code") {
		ctx.Code = params.Get("code")
		ctx.State = params.Get("state")
	}

	// Complete Authentication
	if ctx.Error == "" && ctx.Code != "" {
		if err := s.authorize(ctx); err != nil {
			ctx.Error = "token_failure"
			ctx.Description = err.Error()
		}
	}

	// Render the response to the user
	data := MustAsset("index.html")
	tmpl := template.Must(template.New("index").Parse(string(data)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, ctx); err != nil {
		s.echan <- err
		return
	}

	// Now that authorization has been completed, shutdown the server
	// Must be done in its own go routine to ensure that the body is written first
	// https://medium.com/@int128/shutdown-http-server-by-endpoint-in-go-2a0e2d7f9b8c
	go func() {
		s.echan <- s.Shutdown()
	}()
}

func (s *Server) authorize(ctx RenderContext) (err error) {
	// Check CSRF Protection
	if s.state != ctx.State {
		return errors.New("csrf protection failed: state did not match expected value")
	}

	// Request a token from the Auth0 API
	u := &url.URL{
		Scheme: "https",
		Host:   s.conf.Domain,
		Path:   "/oauth/token",
	}

	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", s.conf.ClientID)
	params.Set("code_verifier", s.codeVerifier)
	params.Set("code", ctx.Code)
	params.Set("redirect_uri", RedirectURI)

	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, u.String(), strings.NewReader(params.Encode())); err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var rep *http.Response
	if rep, err = http.DefaultClient.Do(req); err != nil {
		return err
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		// Attempt to parse JSON error if any
		fmt.Println("An error occurred while fetching the access token:")
		io.Copy(os.Stdout, rep.Body)
		fmt.Println("")

		return fmt.Errorf("could not complete request server responded %s", rep.Status)
	}

	// Write the JSON response to the token cache
	var f *os.File
	if f, err = os.OpenFile(s.conf.TokenCache, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return fmt.Errorf("could not open token cache file: %s", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, rep.Body); err != nil {
		return fmt.Errorf("could not write JSON response into token cache file: %s", err)
	}
	return nil
}

func (s *Server) Favicon(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(MustAsset("favicon.png"))
}

func (s *Server) GetAuthenticationURL() (u *url.URL, err error) {
	u = &url.URL{
		Scheme: "https",
		Host:   s.conf.Domain,
		Path:   "/authorize",
	}

	if s.codeVerifier, err = GenerateCodeToken(); err != nil {
		return nil, err
	}

	if s.state, err = GenerateCodeToken(); err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("audience", s.conf.Audience)
	query.Set("scope", "openid profile email")
	query.Set("response_type", "code")
	query.Set("client_id", s.conf.ClientID)
	query.Set("redirect_uri", RedirectURI)
	query.Set("code_challenge", CodeChallenge(s.codeVerifier))
	query.Set("code_challenge_method", "S256")
	query.Set("state", s.state)

	u.RawQuery = query.Encode()
	return u, nil
}

func GenerateCodeToken() (_ string, err error) {
	nonce := make([]byte, 32)
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(nonce), nil
}

func CodeChallenge(verifier string) string {
	sig := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sig[:])
}

type RenderContext struct {
	Code        string `json:"code"`
	State       string `json:"state"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}
