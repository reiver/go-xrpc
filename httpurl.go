package xrpc

import (
	"strings"

	"github.com/reiver/go-erorr"
)

func httpURL(url string) (string, error) {

	var httpurl string
	{
		var str string = strings.ToLower(url)

		switch {
		case strings.HasPrefix(str, "https://") || strings.HasPrefix(str, "http://"):
			httpurl = url

		case strings.HasPrefix(str, Scheme+"://") || strings.HasPrefix(str, SchemeUnencrypted+"://"):
			xrpcURL, err := ParseURL(url)
			if nil != err {
				return "", err
			}

			httpurl, err = xrpcURL.Resolve(RequestTypeQuery)
			if nil != err {
				return "", err
			}
		default:
			return "", erorr.Errorf("xrpc: an XRPC 'query' can only be made on \"https\", \"http\", %q, and %q not %q", Scheme, SchemeUnencrypted, url)
		}
	}

	return httpurl, nil
}
