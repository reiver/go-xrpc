package xrpc

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrRequest          = erorr.Error("xrpc: request failed")
	ErrResponse         = erorr.Error("xrpc: bad response")
	ErrResponseTooLarge = erorr.Error("xrpc: response too large")
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
