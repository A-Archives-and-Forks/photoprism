package thumb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	t.Run("ValidJPEG", func(t *testing.T) {
		assert.NoError(t, Verify("testdata/example.jpg"))
	})
	t.Run("ValidPNG", func(t *testing.T) {
		assert.NoError(t, Verify("testdata/example.png"))
	})
	t.Run("EmptyFilename", func(t *testing.T) {
		assert.Error(t, Verify(""))
	})
	t.Run("NotAnImage", func(t *testing.T) {
		name := filepath.Join(t.TempDir(), "garbage.jpg")
		assert.NoError(t, os.WriteFile(name, []byte("this is not a valid jpeg payload, only plain text"), 0o644))
		assert.Error(t, Verify(name))
	})
	t.Run("MissingFile", func(t *testing.T) {
		assert.Error(t, Verify(filepath.Join(t.TempDir(), "does-not-exist.jpg")))
	})
}
