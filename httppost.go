package xrpc

import (
	"io"
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-errhttp"
)

func httpPOST(bearerToken string, httpURL string, body io.ReadCloser) (io.ReadCloser, error) {

	var req *http.Request
	{
		var err error
		req, err = http.NewRequest(http.MethodPost, httpURL, nil)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem creating HTTP POST request: %w", err)
		}
	}

	{
		req.Body = body
	}

	{
		err := setUserAgent(req, useragent)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"User-Agent\" header in HTTP POST request: %w", err)
		}
	}

	{
		err := setContentType(req, jsonContentType)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"Content-Type\" header in HTTP POST request: %w", err)
		}
	}

	if 0 < len(bearerToken) {
		err := setAuthorizationBearer(req, bearerToken)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"Authorization\" header in HTTP POST request: %w", err)
		}
	}

	{
		err := setAccept(req, jsonContentType)
		if nil != err {
			return nil, erorr.Errorf("xrpc: problem setting \"Accept\" header in HTTP POST request: %w", err)
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
