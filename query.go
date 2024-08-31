package xrpc

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
	if nil == dst {
		return errNilDestination
	}

	httpurl, err := httpURL(url)
	if nil != err {
		return err
	}

	bodyReadCloser, err := httpGET(httpurl)
	if nil != err {
		return err
	}
	if nil == bodyReadCloser {
		return errNilHTTPResponseBody
	}
	defer bodyReadCloser.Close()

	err = unmarshal(dst, bodyReadCloser)
	if nil != err {
		return err
	}

	return nil
}
