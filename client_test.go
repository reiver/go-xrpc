package xrpc_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/reiver/go-errhttp"

	"github.com/reiver/go-xrpc"
)

// recorded is what the test server saw.
type recorded struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

func newServer(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *[]recorded) {
	t.Helper()

	var requests []recorded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests = append(requests, recorded{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Header: r.Header.Clone(),
			Body:   string(body),
		})
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func TestClient_QueryJSON(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"did":"did:plc:ewvi7nxzyoun6zhxrhs64oiz","handle":"atproto.com"}`)
	})

	var header http.Header = http.Header{}
	header.Set("Authorization", "Bearer abc123")
	header.Set("Atproto-Proxy", "did:web:api.bsky.app#bsky_appview")

	var response struct {
		DID    string `json:"did"`
		Handle string `json:"handle"`
	}
	err := xrpc.Client{UserAgent: "test/1.0"}.QueryJSON(context.Background(), &response, server.URL+"/xrpc/com.atproto.repo.describeRepo?repo=atproto.com", header)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}

	if "did:plc:ewvi7nxzyoun6zhxrhs64oiz" != response.DID || "atproto.com" != response.Handle {
		t.Errorf("Did not get the expected response: %#v", response)
	}

	if 1 != len(*requests) {
		t.Fatalf("Expected 1 request but actually got %d.", len(*requests))
	}
	var request recorded = (*requests)[0]
	if http.MethodGet != request.Method {
		t.Errorf("Expected the method %q but actually got %q.", http.MethodGet, request.Method)
	}
	if "/xrpc/com.atproto.repo.describeRepo" != request.Path || "repo=atproto.com" != request.Query {
		t.Errorf("Did not get the expected path and query: %q %q", request.Path, request.Query)
	}
	for name, expected := range map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer abc123",
		"Atproto-Proxy": "did:web:api.bsky.app#bsky_appview",
		"User-Agent":    "test/1.0",
	} {
		if actual := request.Header.Get(name); expected != actual {
			t.Errorf("Expected the %q header to be %q but actually got %q.", name, expected, actual)
		}
	}
	if "" != request.Header.Get("Content-Type") {
		t.Errorf("Did not expect a \"Content-Type\" header on a query, but actually got %q.", request.Header.Get("Content-Type"))
	}

	// The caller's header is not changed.
	if "" != header.Get("Accept") {
		t.Errorf("Expected the caller's header to be left unchanged, but it got an \"Accept\" header.")
	}
}

func TestClient_Query_raw(t *testing.T) {
	var body string = "\x2a\x6c\x01\x00 not JSON"

	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.dasl.star")
		w.Header().Set("Atproto-Space-Commit", `rev="3l6abc"`)
		io.WriteString(w, body)
	})

	var header http.Header = http.Header{}
	header.Set("Accept", "application/vnd.dasl.star")
	header.Set("Accept-Signature", `mac=("atproto-space-commit";sf);tag="atproto-space-commit"`)

	response, err := xrpc.Client{}.Query(context.Background(), server.URL+"/xrpc/com.atproto.space.getRepo", header)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	if body != string(data) {
		t.Errorf("Expected the body %q but actually got %q.", body, data)
	}
	if expected, actual := `rev="3l6abc"`, response.Header.Get("Atproto-Space-Commit"); expected != actual {
		t.Errorf("Expected the response header %q but actually got %q.", expected, actual)
	}

	var request recorded = (*requests)[0]
	if expected, actual := "application/vnd.dasl.star", request.Header.Get("Accept"); expected != actual {
		t.Errorf("Expected the \"Accept\" header %q but actually got %q.", expected, actual)
	}
	if expected, actual := `mac=("atproto-space-commit";sf);tag="atproto-space-commit"`, request.Header.Get("Accept-Signature"); expected != actual {
		t.Errorf("Expected the \"Accept-Signature\" header %q but actually got %q.", expected, actual)
	}
	if expected, actual := "go-xrpc/0.0 (+https://github.com/reiver/go-xrpc)", request.Header.Get("User-Agent"); expected != actual {
		t.Errorf("Expected the default \"User-Agent\" %q but actually got %q.", expected, actual)
	}
}

func TestClient_ExecuteJSON(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"credential":"eyJ.x.y"}`)
	})

	var response map[string]any
	src := map[string]any{"space": "at://did:plc:ewvi7nxzyoun6zhxrhs64oiz/space/app.example.group/self"}
	err := xrpc.Client{}.ExecuteJSON(context.Background(), &response, server.URL+"/xrpc/com.atproto.space.authorizeSpaceSync", nil, src)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	if "eyJ.x.y" != response["credential"] {
		t.Errorf("Did not get the expected response: %#v", response)
	}

	var request recorded = (*requests)[0]
	if http.MethodPost != request.Method {
		t.Errorf("Expected the method %q but actually got %q.", http.MethodPost, request.Method)
	}
	if expected := `{"space":"at://did:plc:ewvi7nxzyoun6zhxrhs64oiz/space/app.example.group/self"}`; expected != strings.TrimSpace(request.Body) {
		t.Errorf("Expected the body %q but actually got %q.", expected, request.Body)
	}
	if expected, actual := "application/json", request.Header.Get("Content-Type"); expected != actual {
		t.Errorf("Expected the \"Content-Type\" header %q but actually got %q.", expected, actual)
	}
}

func TestClient_ExecuteJSON_noInputNoOutput(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := xrpc.Client{}.ExecuteJSON(context.Background(), nil, server.URL+"/xrpc/com.example.doSomething", nil, nil)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}

	var request recorded = (*requests)[0]
	if "" != request.Body {
		t.Errorf("Expected no body but actually got %q.", request.Body)
	}
	if "" != request.Header.Get("Content-Type") {
		t.Errorf("Did not expect a \"Content-Type\" header with no body, but actually got %q.", request.Header.Get("Content-Type"))
	}
}

func TestClient_Execute_raw(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	response, err := xrpc.Client{}.Execute(context.Background(), server.URL+"/xrpc/com.atproto.repo.uploadBlob", nil, "image/png", strings.NewReader("\x89PNG"))
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	response.Body.Close()

	var request recorded = (*requests)[0]
	if "\x89PNG" != request.Body {
		t.Errorf("Expected the body %q but actually got %q.", "\x89PNG", request.Body)
	}
	if expected, actual := "image/png", request.Header.Get("Content-Type"); expected != actual {
		t.Errorf("Expected the \"Content-Type\" header %q but actually got %q.", expected, actual)
	}
	if "" != request.Header.Get("Accept") {
		t.Errorf("Did not expect an \"Accept\" header, but actually got %q.", request.Header.Get("Accept"))
	}
}

func TestClient_error(t *testing.T) {
	tests := []struct {
		Name            string
		Status          int
		Body            string
		ExpectedName    string
		ExpectedMessage string
		ExpectedHTTP    error
		ExpectedString  string
	}{
		{
			Name:            "XRPC error",
			Status:          http.StatusBadRequest,
			Body:            `{"error":"RepoNotFound","message":"Could not find repo"}`,
			ExpectedName:    "RepoNotFound",
			ExpectedMessage: "Could not find repo",
			ExpectedHTTP:    errhttp.ErrBadRequest,
			ExpectedString:  "xrpc: HTTP 400 RepoNotFound: Could not find repo",
		},
		{
			Name:           "no message",
			Status:         http.StatusUnauthorized,
			Body:           `{"error":"AuthRequired"}`,
			ExpectedName:   "AuthRequired",
			ExpectedHTTP:   errhttp.ErrUnauthorized,
			ExpectedString: "xrpc: HTTP 401 AuthRequired",
		},
		{
			Name:           "not JSON",
			Status:         http.StatusBadGateway,
			Body:           `<html>Bad Gateway</html>`,
			ExpectedHTTP:   errhttp.ErrBadGateway,
			ExpectedString: "xrpc: HTTP 502",
		},
		{
			Name:           "wrong JSON types",
			Status:         http.StatusNotFound,
			Body:           `{"error":123,"message":["x"]}`,
			ExpectedHTTP:   errhttp.ErrNotFound,
			ExpectedString: "xrpc: HTTP 404",
		},
		{
			Name:           "redirect not followed",
			Status:         http.StatusFound,
			Body:           ``,
			ExpectedString: "xrpc: HTTP 302",
		},
	}

	for _, test := range tests {
		server, _ := newServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "30")
			w.Header().Set("Location", "/elsewhere")
			w.WriteHeader(test.Status)
			io.WriteString(w, test.Body)
		})

		var client xrpc.Client = xrpc.Client{
			HTTPClient: &http.Client{
				CheckRedirect: func(*http.Request, []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		}

		var dst map[string]any
		err := client.QueryJSON(context.Background(), &dst, server.URL+"/xrpc/com.atproto.sync.getLatestCommit?did=did:plc:ewvi7nxzyoun6zhxrhs64oiz", nil)

		var xrpcError *xrpc.Error
		if !errors.As(err, &xrpcError) {
			t.Errorf("For %s, expected an *xrpc.Error but actually got: (%T) %v", test.Name, err, err)
			continue
		}
		if test.Status != xrpcError.StatusCode {
			t.Errorf("For %s, expected the status %d but actually got %d.", test.Name, test.Status, xrpcError.StatusCode)
		}
		if test.ExpectedName != xrpcError.Name || test.ExpectedName != xrpc.ErrorName(err) {
			t.Errorf("For %s, expected the name %q but actually got %q (ErrorName: %q).", test.Name, test.ExpectedName, xrpcError.Name, xrpc.ErrorName(err))
		}
		if test.ExpectedMessage != xrpcError.Message {
			t.Errorf("For %s, expected the message %q but actually got %q.", test.Name, test.ExpectedMessage, xrpcError.Message)
		}
		if "30" != xrpcError.Header.Get("Retry-After") {
			t.Errorf("For %s, expected the response headers, but \"Retry-After\" was %q.", test.Name, xrpcError.Header.Get("Retry-After"))
		}
		if nil != test.ExpectedHTTP && !errors.Is(err, test.ExpectedHTTP) {
			t.Errorf("For %s, expected the error to wrap %v, but it did not.", test.Name, test.ExpectedHTTP)
		}
		if test.ExpectedString != err.Error() {
			t.Errorf("For %s, expected the error string %q but actually got %q.", test.Name, test.ExpectedString, err.Error())
		}
	}
}

func TestClient_error_raw(t *testing.T) {
	server, _ := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, `{"error":"InvalidSpaceMac","message":"bad MAC"}`)
	})

	response, err := xrpc.Client{}.Query(context.Background(), server.URL+"/xrpc/com.atproto.space.getRecord", nil)
	if nil != response {
		t.Errorf("Expected no response with an error, but actually got one.")
	}
	if "InvalidSpaceMac" != xrpc.ErrorName(err) {
		t.Errorf("Expected the XRPC error %q but actually got: %v", "InvalidSpaceMac", err)
	}
}

func TestErrorName_notXRPC(t *testing.T) {
	if name := xrpc.ErrorName(errors.New("something else")); "" != name {
		t.Errorf("Expected \"\" but actually got %q.", name)
	}
	if name := xrpc.ErrorName(nil); "" != name {
		t.Errorf("Expected \"\" but actually got %q.", name)
	}
}

func TestClient_MaxResponseSize(t *testing.T) {
	const body = `{"cursor":"abc"}`

	server, _ := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, body)
	})

	var dst map[string]any
	if err := (xrpc.Client{MaxResponseSize: int64(len(body))}).QueryJSON(context.Background(), &dst, server.URL+"/xrpc/com.example.list", nil); nil != err {
		t.Errorf("Did not expect an error at exactly the maximum size, but actually got one: (%T) %s", err, err)
	}

	err := xrpc.Client{MaxResponseSize: int64(len(body) - 1)}.QueryJSON(context.Background(), &dst, server.URL+"/xrpc/com.example.list", nil)
	if !errors.Is(err, xrpc.ErrResponseTooLarge) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", xrpc.ErrResponseTooLarge, err)
	}
}

func TestClient_context(t *testing.T) {
	server, _ := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var dst map[string]any
	err := xrpc.Client{}.QueryJSON(ctx, &dst, server.URL+"/xrpc/com.example.slow", nil)
	if !errors.Is(err, xrpc.ErrRequest) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", xrpc.ErrRequest, err)
	}
}

func TestClient_HTTPClient(t *testing.T) {
	var used bool
	var client xrpc.Client = xrpc.Client{
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				used = true
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Request:    r,
				}, nil
			}),
		},
	}

	var dst map[string]any
	if err := client.QueryJSON(context.Background(), &dst, "xrpc://pds.example.com/com.example.thing", nil); nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	if !used {
		t.Errorf("Expected the client's HTTP client to be used, but it was not.")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func TestAuthorizedQuery(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"ok":true}`)
	})

	var dst map[string]any
	if err := xrpc.AuthorizedQuery(&dst, "tkn", server.URL+"/xrpc/com.example.thing"); nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	if true != dst["ok"] {
		t.Errorf("Did not get the expected response: %#v", dst)
	}
	if expected, actual := "Bearer tkn", (*requests)[0].Header.Get("Authorization"); expected != actual {
		t.Errorf("Expected the \"Authorization\" header %q but actually got %q.", expected, actual)
	}

	if err := xrpc.Query(&dst, server.URL+"/xrpc/com.example.thing"); nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	if actual := (*requests)[1].Header.Get("Authorization"); "" != actual {
		t.Errorf("Did not expect an \"Authorization\" header, but actually got %q.", actual)
	}
}

func TestAuthorizedExecute(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":"RecordNotFound"}`)
	})

	var dst map[string]any
	err := xrpc.AuthorizedExecute(&dst, "tkn", server.URL+"/xrpc/com.example.thing", map[string]string{"a": "b"})
	if !errors.Is(err, errhttp.ErrNotFound) {
		t.Errorf("Expected an error wrapping errhttp.ErrNotFound (as before) but actually got: %v", err)
	}
	if "RecordNotFound" != xrpc.ErrorName(err) {
		t.Errorf("Expected the XRPC error %q but actually got: %v", "RecordNotFound", err)
	}
	if expected, actual := "Bearer tkn", (*requests)[0].Header.Get("Authorization"); expected != actual {
		t.Errorf("Expected the \"Authorization\" header %q but actually got %q.", expected, actual)
	}
}

func TestClient_callerHeadersWin(t *testing.T) {
	server, requests := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{}`)
	})

	var header http.Header = http.Header{}
	header.Set("Accept", "application/json, text/plain")
	header.Set("User-Agent", "caller/2.0")

	var dst map[string]any
	if err := (xrpc.Client{UserAgent: "client/1.0"}).QueryJSON(context.Background(), &dst, server.URL+"/xrpc/com.example.thing", header); nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}

	var request recorded = (*requests)[0]
	if expected, actual := "application/json, text/plain", request.Header.Get("Accept"); expected != actual {
		t.Errorf("Expected the caller's \"Accept\" header %q but actually got %q.", expected, actual)
	}
	if expected, actual := []string{"caller/2.0"}, request.Header.Values("User-Agent"); 1 != len(actual) || expected[0] != actual[0] {
		t.Errorf("Expected the caller's \"User-Agent\" header %q but actually got %q.", expected, actual)
	}
}
