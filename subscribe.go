package xrpc

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/reiver/go-erorr"
)

// Subscribe makes a 'subscribe' XRPC request for the provided XRPC URL in 'url'.
//
// Example usage:
//
//	var response map[string]any = map[string]any{}
//	
//	url = "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"
//	
//	err := xrpc.Subscribe(&response, url)
func Subscribe(url string) (Iterator, error) {
	xrpcURL, err := ParseURL(url)
	if nil != err {
		return nil, erorr.Errorf("xrpc: problem parsing XRPC URL %q: %q", url, err)
	}

	wsURL, err := xrpcURL.resolveWebSocket()
	if nil != err {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{})
	if nil != err {
		return nil, erorr.Errorf("xrpc: problem dialing to websocket at %q (i.e., %q): %w", wsURL, url, err)
	}
	if nil == conn {
		return nil, errNilWebSocketConnection
	}

	var iter Iterator = &internalIterator{conn:conn}

	return iter, nil
}
