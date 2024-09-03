package xrpc

import (
	"io"
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
)

func httpGET(bearerToken string, httpURL string) (io.ReadCloser, error) {

	var req *http.Request
	{
		var err error
		req, err = http.NewRequest(http.MethodGet, httpURL, nil)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem creating HTTP GET request: %w", err)
		}
	}

	{
		err := setUserAgent(req, useragent)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"User-Agent\" header in HTTP GET request: %w", err)
		}
	}

	if 0 < len(bearerToken) {
		err := setAuthorizationBearer(req, bearerToken)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"Authorization\" header in HTTP GET request: %w", err)
		}
	}

	{
		err := setAccept(req, jsonContentType)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"Accept\" header in HTTP GET request: %w", err)
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
