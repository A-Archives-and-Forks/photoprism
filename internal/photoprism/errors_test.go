package photoprism

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestErrCanceled(t *testing.T) {
	assert.EqualError(t, ErrCanceled, "canceled")

	wrapped := fmt.Errorf("worker stopped: %w", ErrCanceled)
	assert.True(t, errors.Is(wrapped, ErrCanceled))
	assert.False(t, errors.Is(wrapped, ErrInsufficientStorage))
}

func TestErrInsufficientStorage(t *testing.T) {
	assert.EqualError(t, ErrInsufficientStorage, "insufficient storage")

	wrapped := fmt.Errorf("walk aborted: %w", ErrInsufficientStorage)
	assert.True(t, errors.Is(wrapped, ErrInsufficientStorage))
	assert.False(t, errors.Is(wrapped, ErrCanceled))
}

func TestWalkResultLog(t *testing.T) {
	t.Run("NilError", func(t *testing.T) {
		level, msg, emit := walkResultLog("index", nil)
		assert.False(t, emit)
		assert.Equal(t, logrus.Level(0), level)
		assert.Empty(t, msg)
	})
	t.Run("Canceled", func(t *testing.T) {
		level, msg, emit := walkResultLog("index", ErrCanceled)
		assert.True(t, emit)
		assert.Equal(t, logrus.InfoLevel, level)
		assert.Equal(t, "index: canceled", msg)
	})
	t.Run("WrappedCanceled", func(t *testing.T) {
		level, _, emit := walkResultLog("index", fmt.Errorf("worker: %w", ErrCanceled))
		assert.True(t, emit)
		assert.Equal(t, logrus.InfoLevel, level)
	})
	t.Run("InsufficientStorage", func(t *testing.T) {
		// Suppressed because the storage helper already emitted the actionable line.
		_, _, emit := walkResultLog("import", ErrInsufficientStorage)
		assert.False(t, emit)
	})
	t.Run("WrappedInsufficientStorage", func(t *testing.T) {
		_, _, emit := walkResultLog("import", fmt.Errorf("walk: %w", ErrInsufficientStorage))
		assert.False(t, emit)
	})
	t.Run("UnknownError", func(t *testing.T) {
		level, msg, emit := walkResultLog("import", errors.New("disk vanished"))
		assert.True(t, emit)
		assert.Equal(t, logrus.ErrorLevel, level)
		assert.Equal(t, "import: disk vanished", msg)
	})
}
