package xrpc

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// Client makes XRPC requests.
//
// The zero value is ready to use.
//
// Example usage:
//
//	var client = xrpc.Client{
//		HTTPClient: httpClient,
//	}
//
//	var header = http.Header{}
//	header.Set("Authorization", "Bearer "+accessToken)
//
//	var response map[string]any
//	err := client.QueryJSON(ctx, &response, "https://pds.example.com/xrpc/com.atproto.repo.describeRepo?repo=did:plc:ewvi7nxzyoun6zhxrhs64oiz", header)
type Client struct {
	// HTTPClient makes the HTTP requests. If nil, http.DefaultClient is used.
	HTTPClient *http.Client

	// UserAgent is sent as the "User-Agent" header (unless the request's header has one). If "", the package's (see [SetUserAgent]) is used.
	UserAgent string

	// MaxResponseSize is the largest response body (in bytes) [Client.QueryJSON] and [Client.ExecuteJSON] read. If 0 (or less), there is no limit.
	//
	// ([Client.Query] and [Client.Execute] return the body unread, for the caller to limit.)
	MaxResponseSize int64
}

// Query makes a 'query' XRPC request (an HTTP GET) to 'url', with the headers in 'header' (which may be nil).
//
// It returns the HTTP response, with its body unread — the caller must close it. (Useful for streamed and non-JSON responses, and for reading response headers.)
//
// If the HTTP status is not 2xx, it returns an [*Error] instead (and closes the body).
//
// 'url' may be an https, http, xrpc, or xrpc-unencrypted URL.
func (receiver Client) Query(ctx context.Context, url string, header http.Header) (*http.Response, error) {
	return receiver.do(ctx, http.MethodGet, url, requestTypeQuery, header, "", nil)
}

// Execute makes an 'execute' XRPC request (an HTTP POST) to 'url', with the headers in 'header' (which may be nil), and with 'body' (which may be nil) as the request body, of type 'contentType'.
//
// NOTE that the official ATProto docs call the 'execute' XRPC request-type 'procedure'.
//
// It returns the HTTP response, with its body unread — the caller must close it.
//
// If the HTTP status is not 2xx, it returns an [*Error] instead (and closes the body).
//
// 'url' may be an https, http, xrpc, or xrpc-unencrypted URL.
func (receiver Client) Execute(ctx context.Context, url string, header http.Header, contentType string, body io.Reader) (*http.Response, error) {
	return receiver.do(ctx, http.MethodPost, url, requestTypeExecute, header, contentType, body)
}

// QueryJSON makes a 'query' XRPC request (see [Client.Query]), and decodes the JSON response into 'dst'.
//
// ('dst' may also be a *string, *[]byte, or *[]rune, to get the response body as is.)
func (receiver Client) QueryJSON(ctx context.Context, dst any, url string, header http.Header) error {
	if nil == dst {
		return errNilDestination
	}

	response, err := receiver.Query(ctx, url, withAccept(header, jsonContentType))
	if nil != err {
		return err
	}
	defer response.Body.Close()

	return receiver.decode(dst, response.Body)
}

// ExecuteJSON makes an 'execute' XRPC request (see [Client.Execute]) whose body is 'src' as JSON (or no body, if 'src' is nil), and decodes the JSON response into 'dst'.
//
// ('dst' may also be a *string, *[]byte, or *[]rune, to get the response body as is. If 'dst' is nil, the response body is ignored — for procedures with no output.)
func (receiver Client) ExecuteJSON(ctx context.Context, dst any, url string, header http.Header, src any) error {
	var body io.Reader
	var contentType string
	if nil != src {
		readCloser, err := marshalJSON(src)
		if nil != err {
			return err
		}
		body = readCloser
		contentType = jsonContentType
	}

	response, err := receiver.Execute(ctx, url, withAccept(header, jsonContentType), contentType, body)
	if nil != err {
		return err
	}
	defer response.Body.Close()

	if nil == dst {
		return nil
	}
	return receiver.decode(dst, response.Body)
}

func (receiver Client) decode(dst any, body io.Reader) error {
	var maxSize int64 = receiver.MaxResponseSize
	if maxSize <= 0 {
		return unmarshalJSON(dst, body)
	}

	data, err := io.ReadAll(io.LimitReader(body, maxSize+1))
	if nil != err {
		return erorr.Wrap(ErrResponse, "problem reading the response body", field.Error("error", err))
	}
	if maxSize < int64(len(data)) {
		return erorr.Wrap(ErrResponseTooLarge, "response body is too large", field.Int64("max", maxSize))
	}

	return unmarshalJSON(dst, bytes.NewReader(data))
}

func (receiver Client) do(ctx context.Context, method string, url string, requestType string, header http.Header, contentType string, body io.Reader) (*http.Response, error) {
	if nil == ctx {
		ctx = context.Background()
	}

	httpURL, err := resolve(url, requestType)
	if nil != err {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, httpURL, body)
	if nil != err {
		return nil, erorr.Wrap(ErrRequest, "problem creating the HTTP request", field.String("method", method), field.String("url", httpURL), field.Error("error", err))
	}

	for name, values := range header {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	if "" != contentType && "" == request.Header.Get("Content-Type") {
		request.Header.Set("Content-Type", contentType)
	}
	if "" == request.Header.Get("User-Agent") {
		var userAgent string = receiver.UserAgent
		if "" == userAgent {
			userAgent = useragent
		}
		if "" != userAgent {
			request.Header.Set("User-Agent", userAgent)
		}
	}

	var client *http.Client = receiver.HTTPClient
	if nil == client {
		client = http.DefaultClient
	}

	response, err := client.Do(request)
	if nil != err {
		return nil, erorr.Wrap(ErrRequest, "HTTP request failed", field.String("method", method), field.String("url", httpURL), field.Error("error", err))
	}
	if nil == response {
		return nil, errNilHTTPResponse
	}
	if nil == response.Body {
		return nil, errNilHTTPResponseBody
	}

	if response.StatusCode < 200 || 299 < response.StatusCode {
		defer response.Body.Close()
		return nil, newError(response)
	}

	return response, nil
}

// withAccept returns a copy of 'header' with an "Accept" header of 'value', unless 'header' already has one.
func withAccept(header http.Header, value string) http.Header {
	if "" != header.Get("Accept") {
		return header
	}

	var result http.Header = header.Clone()
	if nil == result {
		result = http.Header{}
	}
	result.Set("Accept", value)
	return result
}
