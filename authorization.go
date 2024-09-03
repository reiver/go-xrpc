package xrpc

import (
	"net/http"
)

// setAuthorizationBearer sets the HTTP "Authorization: Bearer {TOKEN}" header on an HTTP-request.
//
// Example usage:
//
//	var req *http.Request = ???
//	
//	// ...
//	
//	setAuthorizationBearer(req, "abcde12345")
func setAuthorizationBearer(req *http.Request, value string) error {
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

	value = "Bearer " + value

	header.Add("Authorization", value)
	return nil
}
