package xrpc

import (
	"context"
)

// Query makes a 'query' XRPC request to the provided URL (passed in 'url'), putting the results into 'dst'.
//
// You can provide Query with an `https`, `http`, `xrpc`, or `xrpc-unencrypted` URL.
//
// Here is an example usage of calling Query with an `xrpc` URL:
//
//	var response map[string]any = map[string]any{}
//	
//	url := "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"
//	
//	err := xrpc.Query(&response, url)
//
// Here is an example usage of calling Query with an `https` URL:
//
//	var response map[string]any = map[string]any{}
//	
//	url := "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social"
//	
//	err := xrpc.Query(&response, url)
func Query(dst any, url string) error {
	return AuthorizedQuery(dst, "", url)
}

// AuthorizedQuery is like [Query], but sends 'bearerToken' (if not "") in an "Authorization: Bearer" header.
//
// (For a context, other headers, or another HTTP client, see [Client].)
func AuthorizedQuery(dst any, bearerToken string, url string) error {
	return Client{}.QueryJSON(context.Background(), dst, url, bearerHeader(bearerToken))
}
