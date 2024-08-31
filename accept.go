package xrpc

import (
	"net/http"
)

func setAccept(req *http.Request, value string) error {
	if nil == req {
		return errNilHTTPRequest
	}

	if "" == value {
		return nil
	}

	var header http.Header = req.Header
	{
		if nil == header {
			return errNilHTTPRequestHeader
		}
	}

	header.Add("Accept", value)
	return nil
}
