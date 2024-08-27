package xrpc_test

import (
	"testing"

	"github.com/reiver/go-xrpc"
)

func TestURLResolve(t *testing.T) {

	tests := []struct{
		URL xrpc.URL
		ExpectedQuery string
		ExpectedProcedure string
		ExpectedSubscribe string
	}{
		{
			URL: xrpc.URL{
				Host:              "example.com",
				NSID:                               "net.something.fooBar",
			},
			ExpectedQuery:     "https://example.com/xrpc/net.something.fooBar",
			ExpectedProcedure: "https://example.com/xrpc/net.something.fooBar",
			ExpectedSubscribe:   "wss://example.com/xrpc/net.something.fooBar",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:             "example.com",
				NSID:                              "net.something.fooBar",
			},
			ExpectedQuery:     "http://example.com/xrpc/net.something.fooBar",
			ExpectedProcedure: "http://example.com/xrpc/net.something.fooBar",
			ExpectedSubscribe:   "ws://example.com/xrpc/net.something.fooBar",
		},



		{
			URL: xrpc.URL{
				Host:              "public.api.bsky.app",
				NSID:                                       "app.bsky.actor.getProfile",
			},
			ExpectedQuery:     "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
			ExpectedProcedure: "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
			ExpectedSubscribe:   "wss://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:             "public.api.bsky.app",
				NSID:                                      "app.bsky.actor.getProfile",
			},
			ExpectedQuery:     "http://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
			ExpectedProcedure: "http://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
			ExpectedSubscribe:   "ws://public.api.bsky.app/xrpc/app.bsky.actor.getProfile",
		},



		{
			URL: xrpc.URL{
				Host:              "public.api.bsky.app",
				NSID:                                       "app.bsky.actor.getProfile",
				Query:                                                                "actor=reiver.bsky.social",
			},
			ExpectedQuery:     "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
			ExpectedProcedure: "https://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
			ExpectedSubscribe:  "wss://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
		},
		{
			URL: xrpc.URL{
				Unencrypted: true,
				Host:             "public.api.bsky.app",
				NSID:                                      "app.bsky.actor.getProfile",
				Query:                                                               "actor=reiver.bsky.social",
			},
			ExpectedQuery:     "http://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
			ExpectedProcedure: "http://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
			ExpectedSubscribe:   "ws://public.api.bsky.app/xrpc/app.bsky.actor.getProfile?actor=reiver.bsky.social",
		},
	}

	for testNumber, test := range tests {

		{
			requestType := xrpc.RequestTypeQuery
			expected    :=    test.ExpectedQuery

			actual, err := test.URL.Resolve(requestType)

			if nil != err {
				t.Errorf("For test #%d for XRPC request-type %q, did not expect an error but actually got one.", testNumber, requestType)
				t.Logf("ERROR: (%T) %s", err, err)
				t.Logf("URL: %#v", test.URL)
				continue
			}

			if expected != actual {
				t.Errorf("For test #%d for XRPC request-type %q, the actual resolved-url is not what was expected.", testNumber, requestType)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("URL: %#v", test.URL)
				continue
			}
		}

		{
			requestType := xrpc.RequestTypeProcedure
			expected    :=    test.ExpectedProcedure

			actual, err := test.URL.Resolve(requestType)

			if nil != err {
				t.Errorf("For test #%d for XRPC request-type %q, did not expect an error but actually got one.", testNumber, requestType)
				t.Logf("ERROR: (%T) %s", err, err)
				t.Logf("URL: %#v", test.URL)
				continue
			}

			if expected != actual {
				t.Errorf("For test #%d for XRPC request-type %q, the actual resolved-url is not what was expected.", testNumber, requestType)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("URL: %#v", test.URL)
				continue
			}
		}

		{
			requestType := xrpc.RequestTypeSubscribe
			expected    :=    test.ExpectedSubscribe

			actual, err := test.URL.Resolve(requestType)

			if nil != err {
				t.Errorf("For test #%d for XRPC request-type %q, did not expect an error but actually got one.", testNumber, requestType)
				t.Logf("ERROR: (%T) %s", err, err)
				t.Logf("URL: %#v", test.URL)
				continue
			}

			if expected != actual {
				t.Errorf("For test #%d for XRPC request-type %q, the actual resolved-url is not what was expected.", testNumber, requestType)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("URL: %#v", test.URL)
				continue
			}
		}
	}
}
