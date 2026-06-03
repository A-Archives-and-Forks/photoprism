package event

import "runtime/debug"

// LogPanic logs a recovered panic value together with the current stack trace,
// so that crashes which would otherwise terminate the process without any output
// can be diagnosed from the regular log. It does not stop the process; the caller
// decides whether to exit.
func LogPanic(r any) {
	if r == nil || Log == nil {
		return
	}

	Log.Errorf("panic: %v", r)
	Log.Errorf("stack trace:\n%s", debug.Stack())
}
