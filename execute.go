package xrpc

import (
	"io"
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

func AuthorizedExecute(dst any, bearerToken string, url string, src any) error {
	if nil == dst {
		return errNilDestination
	}

	httpurl, err := httpURL(url)
	if nil != err {
		return err
	}

	var requestReadCloser io.ReadCloser
	{
		if nil != src {
			requestReadCloser, err = marshalJSON(src)
			if nil != err {
				return err
			}
		}
	}

	responseBodyReadCloser, err := httpPOST(bearerToken, httpurl, requestReadCloser)
	if nil != err {
		return err
	}
	if nil == responseBodyReadCloser {
		return errNilHTTPResponseBody
	}
	defer responseBodyReadCloser.Close()

	err = unmarshalJSON(dst, responseBodyReadCloser)
	if nil != err {
		return err
	}

	return nil
}
