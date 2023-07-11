package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

func (s *Server) Login(c *gin.Context) {
	in := &sectigo.AuthenticationRequest{}
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	// Very basic auth as this is only used for staging
	if in.Username != s.conf.Auth.Username || in.Password != s.conf.Auth.Password {
		c.JSON(http.StatusForbidden, Err("could not authenticate user with password"))
		return
	}

	var err error
	out := &sectigo.AuthenticationReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.SignedTokenPair(); err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) Refresh(c *gin.Context) {
	token, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	if _, err = s.tokens.Verify(string(token)); err != nil {
		c.JSON(http.StatusUnauthorized, Err(err))
		return
	}

	out := &sectigo.AuthenticationReply{}
	if out.AccessToken, out.RefreshToken, err = s.tokens.SignedTokenPair(); err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) Authenticate(c *gin.Context) {
	parts := strings.Split(c.GetHeader("Authorization"), " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusForbidden, Err("authentication required"))
		return
	}

	if _, err := s.tokens.Verify(strings.TrimSpace(parts[1])); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Err("invalid credentials"))
		return
	}

	c.Next()
}

// Token time constraint constants.
const (
	accessTokenDuration  = 1 * time.Hour
	refreshTokenDuration = 2 * time.Hour
	accessRefreshOverlap = -1 * time.Hour
)

// Global variables that should really not be changed except between major versions.
var (
	signingMethod = jwt.SigningMethodHS512
)

type Claims struct {
	jwt.RegisteredClaims
	Scopes     []string `json:"scopes,omitempty"`
	FirstLogin bool     `json:"first-login"`
}

// A simple token manager that returns jwt.RegisteredClaims with HS512 signatures.
type Tokens struct {
	subject string
	issuer  string
	scopes  []string
	secret  []byte
}

func NewTokens(conf AuthConfig) (*Tokens, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &Tokens{
		subject: conf.Subject,
		issuer:  conf.Issuer,
		scopes:  conf.Scopes,
		secret:  conf.ParseSecret(),
	}, nil
}

// Verify an access or refresh token after parsing and return its claims.
func (tm *Tokens) Verify(tks string) (claims *Claims, err error) {
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(tks, &Claims{}, tm.keyFunc); err != nil {
		return nil, err
	}

	var ok bool
	if claims, ok = token.Claims.(*Claims); ok && token.Valid {
		if !claims.VerifyIssuer(tm.issuer, true) {
			return nil, fmt.Errorf("invalid issuer %q", claims.Issuer)
		}

		return claims, nil
	}

	return nil, fmt.Errorf("could not parse or verify claims from %T", token.Claims)
}

// Sign an access or refresh token and return the token string.
func (tm *Tokens) Sign(token *jwt.Token) (tks string, err error) {
	return token.SignedString(tm.secret)
}

// Create signed token pair - an access and refresh token.
func (tm *Tokens) SignedTokenPair() (accessToken, refreshToken string, err error) {
	var at *jwt.Token
	if at, err = tm.CreateAccessToken(); err != nil {
		return "", "", err
	}

	var rt *jwt.Token
	if rt, err = tm.CreateRefreshToken(at); err != nil {
		return "", "", err
	}

	if accessToken, err = tm.Sign(at); err != nil {
		return "", "", err
	}

	if refreshToken, err = tm.Sign(rt); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// CreateAccessToken from the verified Google credential payload or from an previous
// token if the access token is being reauthorized from previous credentials. Note that
// the returned token only contains the claims and is unsigned.
func (tm *Tokens) CreateAccessToken() (_ *jwt.Token, err error) {
	// Create the claims for the access token, using access token defaults.
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // ID is randomly generated and shared between access and refresh tokens.
			Subject:   tm.subject,
			Issuer:    tm.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
		},
		Scopes:     tm.scopes,
		FirstLogin: true,
	}
	return jwt.NewWithClaims(signingMethod, claims), nil
}

// CreateRefreshToken from the Access token claims with predefined expiration. Note that
// the returned token only contains the claims and is unsigned.
func (tm *Tokens) CreateRefreshToken(accessToken *jwt.Token) (refreshToken *jwt.Token, err error) {
	accessClaims, ok := accessToken.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not retrieve claims from access token")
	}

	// Create claims for the refresh token from the access token defaults.
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessClaims.ID, // ID is randomly generated and shared between access and refresh tokens.
			Audience:  accessClaims.Audience,
			Issuer:    accessClaims.Issuer,
			Subject:   accessClaims.Subject,
			IssuedAt:  accessClaims.IssuedAt,
			NotBefore: jwt.NewNumericDate(accessClaims.ExpiresAt.Add(accessRefreshOverlap)),
			ExpiresAt: jwt.NewNumericDate(accessClaims.IssuedAt.Add(refreshTokenDuration)),
		},
		Scopes:     accessClaims.Scopes,
		FirstLogin: accessClaims.FirstLogin,
	}

	return jwt.NewWithClaims(signingMethod, claims), nil
}

func (tm *Tokens) keyFunc(token *jwt.Token) (key interface{}, err error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return tm.secret, nil
}
