package i18n

import "strings"

// Response represents an i18n-aware response payload.
//
// Err/Msg carry the server-rendered string in the instance locale (a fallback for non-browser
// consumers); MessageID and MessageParams carry the untranslated source string and its parameters
// so the Web UI can render the message in each user's current UI locale.
type Response struct {
	Code          int    `json:"code"`
	Err           string `json:"error,omitempty"`
	Msg           string `json:"message,omitempty"`
	MessageID     string `json:"messageId,omitempty"`
	MessageParams []any  `json:"messageParams,omitempty"`
	Details       string `json:"details,omitempty"`
}

func (r Response) String() string {
	if r.Err != "" {
		return r.Err
	} else {
		return r.Msg
	}
}

// LowerString returns the lowercased message string.
func (r Response) LowerString() string {
	return strings.ToLower(r.String())
}

func (r Response) Error() string {
	return r.Err
}

// Success reports whether the response code indicates success (2xx).
func (r Response) Success() bool {
	return r.Err == "" && r.Code < 400
}

// NewResponse builds a Response with the given code, message ID, and optional parameters.
// It carries the untranslated source string (ID) and Params so the frontend can render the message
// in the current UI locale, in addition to the server-rendered Err/Msg fallback.
func NewResponse(code int, id Message, params ...any) Response {
	r := Response{Code: code, MessageID: Source(id), MessageParams: params}
	if code < 400 {
		r.Msg = Msg(id, params...)
	} else {
		r.Err = Msg(id, params...)
	}
	return r
}
