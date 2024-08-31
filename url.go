package xrpc

import (
	gourl "net/url"

	"github.com/reiver/go-erorr"
	libnsid "github.com/reiver/go-nsid"
)

// URL represents an 'xrpc' and 'xrpc-unencrypted' URL.
//
// For example:
//
//	xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social
//
//	xrpc://example.com/com.example.fooBar
//
//	xrpc-unencrypted://example.com/com.example.fooBar
type URL struct {
	Unencrypted bool
	Host string
	NSID string
	Query string
}

// ConstructURL constructs an XRPC URL (i.e., "xrpc://...") from a 'host', an 'nsid', and a 'query'.
//
// The 'query' can be an empty string.
func ConstructURL(host string, nsid string, query string) URL {
	return URL{
		Host:host,
		NSID:nsid,
		Query:query,
	}
}

// ConstructURL constructs an XRPC URL (i.e., "xrpc-unencrypted://...") from a 'host', an 'nsid', and a 'query'.
//
// The 'query' can be an empty string.
func ConstructUnencryptedURL(host string, nsid string, query string) URL {
	return URL{
		Unencrypted:true,
		Host:host,
		NSID:nsid,
		Query:query,
	}
}

// MustParseURL is similar to [ParseURL] except that it panic()s if there is an error (rather than returning it, like how [ParseURL] does).
func MustParseURL(url string) URL {
	xrpcurl, err := ParseURL(url)
	if nil != err {
		panic(err)
	}

	return xrpcurl
}

func ParseURL(url string) (URL, error) {
	var empty URL

	urloc, err := gourl.Parse(url)
	if nil != err {
		return empty, err
	}
	if nil == urloc {
		return empty, errNilURL
	}

	var nsid string = urloc.Path
	if 0 < len(nsid) && '/' == nsid[0] {
		nsid = nsid[1:]
	}

	switch urloc.Scheme {
	case Scheme:
		return ConstructURL(urloc.Host, nsid, urloc.RawQuery), nil
	case SchemeUnencrypted:
		return ConstructUnencryptedURL(urloc.Host, nsid, urloc.RawQuery), nil
	default:
		return empty, erorr.Errorf("xrpc: expected scheme to be %q or %q but was %q", Scheme, SchemeUnencrypted, urloc.Scheme)
	}

}

// MustResolve is similar to [Resolve] except that it panic()s if there is an error (rather than returning it, like how [Resolve] does).
func (receiver URL) MustResolve(requestType string) string {
	httpurl, err := receiver.Resolve(requestType)
	if nil != err {
		panic(err)
	}

	return httpurl
}

// Resolve turns an XRPC URL into an HTTP URL.
//
// For example:
//
//	"xrpc://public.api.bsky.app/app.bsky.actor.getProfile" -> "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile"
//
//	"xrpc://example.com/com.example.fooBar" -> "https://example.com/xrpc/com.example.fooBar"
//
//	"xrpc-unencrypted://example.com/com.example.fooBar" -> "http://example.com/xrpc/com.example.fooBar"
func (receiver URL) Resolve(requestType string) (string, error) {
	switch requestType {
	case RequestTypeExecute, RequestTypeProcedure:
		return receiver.resolveWeb()
	case RequestTypeQuery:
		return receiver.resolveWeb()
	case RequestTypeSubscribe:
		return receiver.resolveWebSocket()
	default:
		return "", erorr.Errorf("xrpc: do not know how to resolve a %q XRPC request", requestType)
	}
}

func (receiver URL) resolveWeb() (string, error) {
	if err := receiver.Validate(); nil != err {
		return "", err
	}

	var buffer [256]byte
	var p []byte = buffer[0:0]

	p = append(p, "http"...)
	if !receiver.Unencrypted {
		p = append(p, 's')
	}
	p = append(p, "://"...)

	p = append(p, receiver.Host...)

	p = append(p, "/xrpc/"...)
	p = append(p, receiver.NSID...)

	if "" != receiver.Query {
		p = append(p, '?')
		p = append(p, receiver.Query...)
	}

	return string(p), nil
}


func (receiver URL) resolveWebSocket() (string, error) {
	url, err := receiver.resolveWeb()
	if nil != err {
		return "", err
	}

	url = "ws" + url[len("http"):]

	return url, nil
}

// String returns the (serialized) XRPC URL.
func (receiver URL) String() string {
	var buffer [256]byte
	var p []byte = buffer[0:0]

	if receiver.Unencrypted {
		p = append(p, SchemeUnencrypted...)
	} else {
		p = append(p, Scheme...)
	}
	p = append(p, "://"...)

	p = append(p, receiver.Host...)
	p = append(p, '/')
	p = append(p, receiver.NSID...)

	if "" != receiver.Query {
		p = append(p, '?')
		p = append(p, receiver.Query...)
	}

	return string(p)
}

// Validate returns an error if the URL is invalid.
func (receiver URL) Validate() error {
	if "" == receiver.Host {
		return errEmptyHost
	}

	{
		err := libnsid.Validate(receiver.NSID)
		if nil != err {
			return err
		}
	}

	return nil
}
