package xrpc

import (
	"net/http"
)

// bearerHeader returns the headers with "Authorization: Bearer {TOKEN}" — or no headers, if 'bearerToken' is "".
func bearerHeader(bearerToken string) http.Header {
	if "" == bearerToken {
		return nil
	}

	return http.Header{
		"Authorization": {"Bearer " + bearerToken},
	}
}
