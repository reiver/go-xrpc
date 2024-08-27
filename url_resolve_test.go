package xrpc_test

import (
	"testing"

	"github.com/reiver/go-xrpc"
)

func TestURLResolve(t *testing.T) {

	tests := []struct{
		URL xrpc.URL
		Expected string
	}{
		{
			URL: xrpc.URL{
				Host:     "example.com",
				NSID:                      "net.something.fooBar",
			},
			Expected: "https://example.com/xrpc/net.something.fooBar",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:    "example.com",
				NSID:                     "net.something.fooBar",
			},
			Expected: "http://example.com/xrpc/net.something.fooBar",
		},



		{
			URL: xrpc.URL{
				Host:     "public.api.bsky.app",
				NSID:                              "app.bsky.actor.getProfile",
			},
			Expected: "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:    "public.api.bsky.app",
				NSID:                             "app.bsky.actor.getProfile",
			},
			Expected: "http://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
		},



		{
			URL: xrpc.URL{
				Host:     "public.api.bsky.app",
				NSID:                              "app.bsky.actor.getProfile",
				Query:                                                       "actor=reiver.bsky.social",
			},
			Expected: "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:    "public.api.bsky.app",
				NSID:                             "app.bsky.actor.getProfile",
				Query:                                                      "actor=reiver.bsky.social",
			},
			Expected: "http://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
		},
	}

	for testNumber, test := range tests {

		actual, err := test.URL.Resolve()

		if nil != err {
			t.Errorf("For test #%d, did not expect an error but actually got one.", testNumber)
			t.Logf("ERROR: (%T) %s", err, err)
			t.Logf("URL: %#v", test.URL)
			continue
		}

		{
			expected := test.Expected

			if expected != actual {
				t.Errorf("For test #%d, the actual resolved-url is not what was expected.", testNumber)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("URL: %#v", test.URL)
				continue
			}
		}
	}
}
