package xrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/reiver/go-errhttp"
)

// Query makes a 'query' XRPC request for the provided XRPC URL in 'url', putting the results into 'dst'.
//
// Example usage:
//
//	var response map[string]any = map[string]any{}
//	
//	url := "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"
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

	httpResp, err := http.Get(httpURL)
	if nil != err {
		return err
	}
	if nil == httpResp {
		return errNilHTTPResponse
	}

	if http.StatusOK != httpResp.StatusCode {
		return errhttp.Return(httpResp.StatusCode)
	}

	var bodyRC io.ReadCloser = httpResp.Body
	if nil == bodyRC {
		return errNilHTTPResponseBody
	}

	var bodyBuffer bytes.Buffer
	io.Copy(&bodyBuffer, httpResp.Body)

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
