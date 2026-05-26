package photoprism

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

// ErrCanceled signals that a walk callback or worker was interrupted by Cancel().
var ErrCanceled = errors.New("canceled")

// ErrInsufficientStorage signals that a walk callback aborted because the storage
// path is critically low on free disk space or the configured quota is exhausted.
var ErrInsufficientStorage = errors.New("insufficient storage")

// walkResultLog returns the level, message, and emit flag for a directory-walk result.
// ErrCanceled surfaces at info level (user-initiated stop); ErrInsufficientStorage is
// suppressed because the storage-check helper already logged the actionable cause.
func walkResultLog(prefix string, err error) (level logrus.Level, message string, emit bool) {
	switch {
	case err == nil:
		return 0, "", false
	case errors.Is(err, ErrCanceled):
		return logrus.InfoLevel, fmt.Sprintf("%s: canceled", prefix), true
	case errors.Is(err, ErrInsufficientStorage):
		return 0, "", false
	default:
		return logrus.ErrorLevel, fmt.Sprintf("%s: %s", prefix, err.Error()), true
	}
}

// logWalkResult emits a log line for a directory-walk result via walkResultLog.
func logWalkResult(prefix string, err error) {
	if level, msg, emit := walkResultLog(prefix, err); emit {
		log.Log(level, msg)
	}
}
