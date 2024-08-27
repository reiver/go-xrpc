package xrpc

import (
	"github.com/reiver/go-erorr"
)

const (
	errClosed                 = erorr.Error("xrpc: closed")
	errEmptyHost              = erorr.Error("xrpc: empty host")
	errNilDestination         = erorr.Error("xrpc: nil destination")
	errNilHTTPResponse        = erorr.Error("xrpc: nil http-response")
	errNilHTTPResponseBody    = erorr.Error("xrpc: nil http-response body")
	errNilReceiver            = erorr.Error("xrpc: nil receiver")
	errNilURL                 = erorr.Error("xrpc: nil url")
	errNilWebSocketConnection = erorr.Error("xrpc: nil websocket connection")
)
