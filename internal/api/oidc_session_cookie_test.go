package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/rnd"
)

// findCookie returns the named cookie from the recorder, or nil if absent.
func findCookie(w *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, ck := range w.Result().Cookies() {
		if ck.Name == name {
			return ck
		}
	}
	return nil
}

func TestSetOIDCSessionCookie(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		token := rnd.AuthToken()
		SetOIDCSessionCookie(c, token, 3600, true)
		ck := findCookie(w, OIDCSessionCookie)
		if assert.NotNil(t, ck) {
			assert.Equal(t, token, ck.Value)
			assert.Equal(t, "/oauth", ck.Path)
			assert.True(t, ck.HttpOnly)
			assert.True(t, ck.Secure)
			assert.Equal(t, http.SameSiteLaxMode, ck.SameSite)
			assert.Equal(t, 3600, ck.MaxAge)
		}
	})
	t.Run("InsecureOmitsSecureFlag", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		SetOIDCSessionCookie(c, rnd.AuthToken(), 3600, false)
		ck := findCookie(w, OIDCSessionCookie)
		if assert.NotNil(t, ck) {
			assert.False(t, ck.Secure)
		}
	})
	t.Run("InvalidTokenSetsNothing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		SetOIDCSessionCookie(c, "not-a-valid-token", 3600, true)
		assert.Nil(t, findCookie(w, OIDCSessionCookie))
	})
	t.Run("NilContext", func(t *testing.T) {
		assert.NotPanics(t, func() { SetOIDCSessionCookie(nil, rnd.AuthToken(), 3600, true) })
	})
}

func TestClearOIDCSessionCookie(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ClearOIDCSessionCookie(c, false)
		ck := findCookie(w, OIDCSessionCookie)
		if assert.NotNil(t, ck) {
			assert.Equal(t, "", ck.Value)
			assert.Equal(t, "/oauth", ck.Path)
			assert.True(t, ck.MaxAge < 0)
		}
	})
	t.Run("NilContext", func(t *testing.T) {
		assert.NotPanics(t, func() { ClearOIDCSessionCookie(nil, false) })
	})
}

func TestOIDCSessionCookieToken(t *testing.T) {
	newCtx := func(ck *http.Cookie) *gin.Context {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/oauth/authorize", nil)
		if ck != nil {
			req.AddCookie(ck)
		}
		c.Request = req
		return c
	}
	t.Run("ValidToken", func(t *testing.T) {
		token := rnd.AuthToken()
		c := newCtx(&http.Cookie{Name: OIDCSessionCookie, Value: token})
		assert.Equal(t, token, OIDCSessionCookieToken(c))
	})
	t.Run("MalformedToken", func(t *testing.T) {
		c := newCtx(&http.Cookie{Name: OIDCSessionCookie, Value: "tooshort"})
		assert.Equal(t, "", OIDCSessionCookieToken(c))
	})
	t.Run("Absent", func(t *testing.T) {
		assert.Equal(t, "", OIDCSessionCookieToken(newCtx(nil)))
	})
	t.Run("NilContext", func(t *testing.T) {
		assert.Equal(t, "", OIDCSessionCookieToken(nil))
	})
}
