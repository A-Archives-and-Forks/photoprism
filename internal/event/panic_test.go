package event

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogPanic(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer
		logger := logrus.New()
		logger.SetOutput(&buf)
		logger.SetLevel(logrus.ErrorLevel)
		orig := Log
		Log = logger
		t.Cleanup(func() { Log = orig })

		LogPanic("boom")

		out := buf.String()
		assert.Contains(t, out, "panic: boom")
		assert.Contains(t, out, "stack trace")
		assert.Contains(t, out, "goroutine")
	})
	t.Run("Nil", func(t *testing.T) {
		var buf bytes.Buffer
		logger := logrus.New()
		logger.SetOutput(&buf)
		orig := Log
		Log = logger
		t.Cleanup(func() { Log = orig })

		LogPanic(nil)

		assert.Empty(t, buf.String())
	})
	t.Run("NoLogger", func(t *testing.T) {
		orig := Log
		Log = nil
		t.Cleanup(func() { Log = orig })

		assert.NotPanics(t, func() { LogPanic("boom") })
	})
}
