package xrpc_test

import (
	"testing"

	"github.com/reiver/go-xrpc"
)

func TestURLString(t *testing.T) {

	tests := []struct{
		URL xrpc.URL
		Expected string
	}{
		{
			URL: xrpc.URL{
				Host:     "example.com",
				NSID:                "net.something.fooBar",
			},
			Expected: "xrpc://example.com/net.something.fooBar",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:    "example.com",
				NSID:                            "net.something.fooBar",
			},
			Expected: "xrpc-unencrypted://example.com/net.something.fooBar",
		},



		{
			URL: xrpc.URL{
				Host:     "public.api.bsky.app",
				NSID:                        "app.bsky.actor.getProfile",
			},
			Expected: "xrpc://public.api.bsky.app/app.bsky.actor.getProfile",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:                "public.api.bsky.app",
				NSID:                                    "app.bsky.actor.getProfile",
			},
			Expected: "xrpc-unencrypted://public.api.bsky.app/app.bsky.actor.getProfile",
		},
	}

	for testNumber, test := range tests {

		actual := test.URL.String()

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
