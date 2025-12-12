package server

import (
	"bytes"
	stdgzip "compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoprism/photoprism/internal/config"
)

func TestGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Enable gzip for this test router.
	conf := config.TestConfig()
	conf.Options().HttpCompression = "gzip"

	r := gin.New()
	r.Use(gzip.Gzip(
		gzip.DefaultCompression,
		gzip.WithExcludedExtensions([]string{
			".png", ".gif", ".jpeg", ".jpg", ".webp", ".mp3", ".mp4", ".zip", ".gz",
		}),
		gzip.WithExcludedPaths([]string{
			conf.BaseUri("/health"),
			conf.BaseUri(config.ApiUri + "/t"),
			conf.BaseUri(config.ApiUri + "/folders/t"),
			conf.BaseUri(config.ApiUri + "/dl"),
			conf.BaseUri(config.ApiUri + "/zip"),
			conf.BaseUri(config.ApiUri + "/albums"),
			conf.BaseUri(config.ApiUri + "/labels"),
			conf.BaseUri(config.ApiUri + "/videos"),
		}),
	))

	r.GET("/ok", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	excludedPath := conf.BaseUri(config.ApiUri + "/dl/test")
	r.GET(excludedPath, func(c *gin.Context) {
		c.String(http.StatusOK, "download")
	})

	t.Run("CompressesSuccessfulResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ok", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

		zr, err := stdgzip.NewReader(bytes.NewReader(w.Body.Bytes()))
		require.NoError(t, err)
		defer zr.Close()

		b, err := io.ReadAll(zr)
		require.NoError(t, err)
		assert.Equal(t, "hello world", string(b))
	})

	t.Run("DoesNotCompressExcludedPaths", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", excludedPath, nil)
		req.Header.Set("Accept-Encoding", "gzip")

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Content-Encoding"))
		assert.Equal(t, "download", w.Body.String())
	})

	t.Run("DoesNotCompressNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/missing", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		assert.Empty(t, w.Header().Get("Content-Encoding"))
		assert.Contains(t, w.Body.String(), "404")
	})
}
