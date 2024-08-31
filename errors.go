package xrpc

import (
	"github.com/reiver/go-erorr"
)

const (
	errClosed                 = erorr.Error("xrpc: closed")
	errEmptyHost              = erorr.Error("xrpc: empty host")
	errNilDestination         = erorr.Error("xrpc: nil destination")
	errNilHTTPRequest         = erorr.Error("xrpc: nil http-request")
	errNilHTTPRequestHeader   = erorr.Error("xrpc: nil http-request header")
	errNilHTTPResponse        = erorr.Error("xrpc: nil http-response")
	errNilHTTPResponseBody    = erorr.Error("xrpc: nil http-response body")
	errNilReader              = erorr.Error("xrpc: nil reader")
	errNilReceiver            = erorr.Error("xrpc: nil receiver")
	errNilURL                 = erorr.Error("xrpc: nil url")
	errNilWebSocketConnection = erorr.Error("xrpc: nil websocket connection")
)
