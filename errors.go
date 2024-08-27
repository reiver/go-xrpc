package xrpc

import (
	"github.com/reiver/go-erorr"
)

const (
	errEmptyHost           = erorr.Error("xrpc: empty host")
	errNilDestination      = erorr.Error("xrpc: nil destination")
	errNilHTTPResponse     = erorr.Error("xrpc: nil http-response")
	errNilHTTPResponseBody = erorr.Error("xrpc: nil http-response body")
	errNilURL              = erorr.Error("xrpc: nil url")
)
