package xrpc

import (
	"io"
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
)

func httpGET(httpURL string) (io.ReadCloser, error) {

	var req *http.Request
	{
		var err error
		req, err = http.NewRequest(http.MethodGet, httpURL, nil)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem creating HTTP request: %w", err)
		}
	}

	{
		err := setUserAgent(req, useragent)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"User-Agent\" in HTTP request: %w", err)
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
