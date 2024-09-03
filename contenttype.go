package xrpc

import (
	"net/http"
)

// setContentType sets the HTTP Content-Type header on an HTTP-request.
//
// Example usage:
//
//	var req *http.Request = ???
//	
//	// ...
//	
//	setContentType(req, "application/json")
func setContentType(req *http.Request, value string) error {
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

	header.Add("Content-Type", value)
	return nil
}
