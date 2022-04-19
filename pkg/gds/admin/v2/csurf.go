package admin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Parameters and names for double-cookie submit CSRF protection
const (
	CSRFCookie          = "csrf_token"
	CSRFReferenceCookie = "csrf_reference_token"
	CSRFHeader          = "X-CSRF-TOKEN"
)

// DoubleCookie is Cross Site Request Forgery (CSRF/XSRF) protection middleware that
// checks the presence of an X-CSRF-TOKEN header containing a cryptographically random
// token that matches a token contained in the CSRF-TOKEN cookie in the request.
// Because of the same-origin policy, an attacker cannot access the cookies or scripts
// of the safe site therefore it will not be able to guess what token to put in the
// request, thereby protecting from CSRF attacks. Note that Double Cookie CSRF
// protection requires TLS to prevent MITM attacks.
func DoubleCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(CSRFReferenceCookie)
		if err != nil {
			log.Warn().Err(err).Msg("no csrf token cookie in request")
			c.JSON(http.StatusForbidden, ErrorResponse(ErrCSRFVerification))
			c.Abort()
			return
		}

		header := c.GetHeader(CSRFHeader)
		if header, err = url.QueryUnescape(header); err != nil {
			log.Warn().Err(err).Msg("could not unescape csrf token")
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			c.Abort()
			return
		}

		if cookie != header {
			log.Warn().Bool("header_exists", header != "").Bool("cookie_exists", cookie != "").Msg("csrf token cookie/header mismatch")
			c.JSON(http.StatusForbidden, ErrorResponse(ErrCSRFVerification))
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetDoubleCookieTokens is a helper function to set cookies on a gin-request.
// The exp parameter is the Unix timestamp the cookie should be expired, which in most
// cases is extracted from the exp field of the refresh token claims.
func SetDoubleCookieTokens(c *gin.Context, domain string, exp time.Time) error {
	// Generate the CSRF token
	token, err := GenerateCSRFToken()
	if err != nil {
		return err
	}

	// Compute max age from the expires unix timestamp of the refresh token.
	maxAge := int((time.Until(exp)).Seconds()) + 60

	// Set the reference cookie
	c.SetCookie(CSRFReferenceCookie, token, maxAge, "/", domain, true, true)

	// Set the csrf token cookie
	c.SetCookie(CSRFCookie, token, maxAge, "/", domain, true, false)
	return nil
}

// This is random seed that is meant to present an additional barrier to cryptanalysis
// and is unique to each process in the network.
var seed []byte

func GenerateCSRFToken() (_ string, err error) {
	// If the process seed is not generated, generate it now.
	if len(seed) == 0 {
		seed = make([]byte, 16)
		if _, err = rand.Read(seed); err != nil {
			return "", err
		}
	}

	nonce := make([]byte, 32)
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	sig := sha256.New()
	sig.Write(seed)
	sig.Write(nonce)

	return base64.URLEncoding.EncodeToString(sig.Sum(nil)), nil
}
