package xrpc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/reiver/go-errhttp"
)

// maxErrorBodySize is how much of an error response's body is read, looking for the XRPC error.
const maxErrorBodySize = 64 * 1024

// Error is an XRPC error — the answer to an XRPC request whose HTTP status was not 2xx.
//
// XRPC errors look like this:
//
//	{"error":"RepoNotFound","message":"Could not find repo for DID: did:plc:ewvi7nxzyoun6zhxrhs64oiz"}
//
// Use errors.As to get it:
//
//	var xrpcError *xrpc.Error
//	if errors.As(err, &xrpcError) && "RepoNotFound" == xrpcError.Name {
//		// ...
//	}
//
// It unwraps to the github.com/reiver/go-errhttp error for its HTTP status (ex: errhttp.ErrNotFound).
type Error struct {
	StatusCode int

	// Name is the response body's "error" (ex: "RepoNotFound"). It is "" if the body did not have one.
	Name string

	// Message is the response body's "message". It is "" if the body did not have one.
	Message string

	// Header is the response's HTTP headers (ex: for "Retry-After", or "WWW-Authenticate").
	Header http.Header
}

var _ error = &Error{}

func (receiver *Error) Error() string {
	if nil == receiver {
		return "xrpc: nil error"
	}

	var s string = "xrpc: HTTP " + strconv.Itoa(receiver.StatusCode)
	if "" != receiver.Name {
		s += " " + receiver.Name
	}
	if "" != receiver.Message {
		s += ": " + receiver.Message
	}
	return s
}

func (receiver *Error) Unwrap() error {
	if nil == receiver {
		return nil
	}
	return errhttp.Return(receiver.StatusCode)
}

// ErrorName returns the name of the XRPC error (ex: "RepoNotFound") in 'err' — see [Error].
// It returns "" if 'err' is not (and does not wrap) an XRPC error, or if the XRPC error has no name.
func ErrorName(err error) string {
	var xrpcError *Error
	if !errors.As(err, &xrpcError) {
		return ""
	}
	return xrpcError.Name
}

// newError returns the [Error] for a (non-2xx) HTTP response. It reads (some of) the body, but does not close it.
func newError(response *http.Response) *Error {
	var xrpcError *Error = &Error{
		StatusCode: response.StatusCode,
		Header:     response.Header,
	}

	if nil == response.Body {
		return xrpcError
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, maxErrorBodySize))
	if nil != err {
		return xrpcError
	}

	var body struct {
		Error   any `json:"error"`
		Message any `json:"message"`
	}
	if err := json.Unmarshal(data, &body); nil != err {
		return xrpcError
	}

	if name, ok := body.Error.(string); ok {
		xrpcError.Name = name
	}
	if message, ok := body.Message.(string); ok {
		xrpcError.Message = message
	}

	return xrpcError
}
