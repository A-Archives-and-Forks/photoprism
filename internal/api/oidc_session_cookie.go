package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// OIDCSessionCookie is the name of the narrowly-scoped cookie that lets the
// Portal OIDC OP authenticate a browser on a top-level navigation to
// /oauth/authorize, which carries no Authorization or X-Auth-Token header.
const OIDCSessionCookie = "oidc_session"

// oidcSessionCookiePath scopes the OP session cookie to the OAuth endpoints so
// it is never transmitted to the general API surface.
const oidcSessionCookiePath = "/oauth"

// SetOIDCSessionCookie stores the Portal session token in a narrowly-scoped,
// HttpOnly cookie so the OIDC OP /oauth/authorize endpoint can authenticate the
// browser on a top-level navigation. The cookie is honored ONLY by the OP
// authorize handler (see OIDCSessionCookieToken), never as a general API
// authenticator, so it adds no CSRF surface to state-changing endpoints.
func SetOIDCSessionCookie(c *gin.Context, authToken string, maxAge int, secure bool) {
	if c == nil || !rnd.IsAuthToken(authToken) {
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     OIDCSessionCookie,
		Value:    authToken,
		Path:     oidcSessionCookiePath,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearOIDCSessionCookie removes the OP session cookie, e.g. on logout.
func ClearOIDCSessionCookie(c *gin.Context, secure bool) {
	if c == nil {
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     OIDCSessionCookie,
		Value:    "",
		Path:     oidcSessionCookiePath,
		MaxAge:   -1,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// OIDCSessionCookieToken returns the session token from the OP session cookie,
// or "" if absent or malformed. It is the only reader of OIDCSessionCookie and
// must be used solely by the OIDC OP authorize handler as a fallback when no
// Authorization/X-Auth-Token header is present on a browser navigation.
func OIDCSessionCookieToken(c *gin.Context) string {
	if c == nil {
		return ""
	}

	v, err := c.Cookie(OIDCSessionCookie)
	if err != nil {
		return ""
	}

	if token := clean.Token(v); rnd.IsAuthToken(token) {
		return token
	}

	return ""
}
