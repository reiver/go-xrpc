package xrpc

import (
	"net/http"
)

// useragent holds the default HTTP User-Agent included in HTTP requests.
//
// Its initial value is what is set below.
//
// Its value can be changed by calling xrpc.SetUserAgent()
var useragent string = "go-xrpc/0.0 (+https://github.com/reiver/go-xrpc)"

// SetUserAgent sets the User-Agent in XRPC requests.
func SetUserAgent(value string) {
	useragent = value
}

func setUserAgent(req *http.Request, useragent string) error {
	if nil == req {
		return errNilHTTPRequest
	}

	if "" == useragent {
		return nil
	}

	var header http.Header = req.Header
	{
		if nil == header {
			return errNilHTTPRequestHeader
		}
	}

	header.Add("User-Agent", useragent)
	return nil
}
