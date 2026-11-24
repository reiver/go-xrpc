package xrpc

import (
	"context"
)

// Execute makes an 'execute' XRPC request to the provided URL (passed in 'url'), putting the results into 'dst'.
//
// NOTE that the official Bluesky docs call the 'execute' XRPC request-type 'procedure'.
// ('execute' was chosen here as it is a verb.)
//
// You can provide Execute with an `https`, `http`, `xrpc`, or `xrpc-unencrypted` URL.
//
// Here is an example usage of calling Execute with an `xrpc` URL:
//
//	var response map[string]any = map[string]any{}
//	
//	src := struct{
//		Identifier string `json:"identifier"`
//		Password   string `json:"password"`
//	}{
//		Identifier: "joeblow.bsky.social",
//		Password:   "password123",
//	}
//	
//	url := "xrpc://bsky.social/com.atproto.server.createSession"
//	
//	err := xrpc.Execute(&response, url, src)
//
// Here is an example usage of calling Execute with an `https` URL:
//
//	var response map[string]any = map[string]any{}
//	
//	src := struct{
//		Identifier string `json:"identifier"`
//		Password   string `json:"password"`
//	}{
//		Identifier: "joeblow.bsky.social",
//		Password:   "password123",
//	}
//	
//	url := "https://bsky.social/xrpc/com.atproto.server.createSession"
//	
//	err := xrpc.Execute(&response, url, src)
func Execute(dst any, url string, src any) error {
	return AuthorizedExecute(dst, "", url, src)
}

// AuthorizedExecute is like [Execute], but sends 'bearerToken' (if not "") in an "Authorization: Bearer" header.
//
// (For a context, other headers, or another HTTP client, see [Client].)
func AuthorizedExecute(dst any, bearerToken string, url string, src any) error {
	if nil == dst {
		return errNilDestination
	}

	return Client{}.ExecuteJSON(context.Background(), dst, url, bearerHeader(bearerToken), src)
}
