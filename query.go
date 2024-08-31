package xrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
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
	if nil == dst {
		return errNilDestination
	}

	xrpcURL, err := ParseURL(url)
	if nil != err {
		return err
	}

	httpURL, err := xrpcURL.Resolve(RequestTypeQuery)
	if nil != err {
		return err
	}

	bodyReadCloser, err := query(httpURL)
	if nil != err {
		return err
	}
	if nil == bodyReadCloser {
		return errNilHTTPResponseBody
	}
	defer bodyReadCloser.Close()

	var bodyBuffer bytes.Buffer
	io.Copy(&bodyBuffer, bodyReadCloser)

	var body []byte = bodyBuffer.Bytes()

	switch casted := dst.(type) {
	case json.Unmarshaler:
		return json.Unmarshal(body, casted)
	case *string:
		*casted = string(body)
		return nil
	case *[]byte:
		*casted = body
		return nil
	case *[]rune:
		*casted = []rune(string(body))
		return nil
	default:
		return json.Unmarshal(body, dst)
	}

	return nil
}

func query(httpURL string) (io.ReadCloser, error) {

	var req *http.Request
	{
		var err error
		req, err = http.NewRequest(http.MethodGet, httpURL, nil)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem creating HTTP request: %w", err)
		}
	}

	httpResp, err := http.DefaultClient.Do(req)
	if nil != err {
		return nil, err
	}
	if nil == httpResp {
		return nil, errNilHTTPResponse
	}

	if http.StatusOK != httpResp.StatusCode {
		return nil, errhttp.Return(httpResp.StatusCode)
	}

	var bodyRC io.ReadCloser = httpResp.Body
	if nil == bodyRC {
		return nil, errNilHTTPResponseBody
	}

	return bodyRC, nil
}
