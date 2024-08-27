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
//	xrpc://public.api.bsky.app/app.bsky.actor.getProfile
//
//	xrpc://example.com/com.example.fooBar
//
//	xrpc-unencrypted://example.com/com.example.fooBar
type URL struct {
	Unencrypted bool
	Host string
	NSID string
}

func ConstructURL(host string, nsid string) URL {
	return URL{
		Host:host,
		NSID:nsid,
	}
}

func ConstructUnencryptedURL(host string, nsid string) URL {
	return URL{
		Unencrypted:true,
		Host:host,
		NSID:nsid,
	}
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
		return ConstructURL(urloc.Host, nsid), nil
	case SchemeUnencrypted:
		return ConstructUnencryptedURL(urloc.Host, nsid), nil
	default:
		return empty, erorr.Errorf("xrpc: expected scheme to be %q or %q but was %q", Scheme, SchemeUnencrypted, urloc.Scheme)
	}

}

func (receiver URL) Resolve() (string, error) {
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

	return string(p), nil
}

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

	return string(p)
}

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
