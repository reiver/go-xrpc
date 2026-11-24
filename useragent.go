package xrpc

// useragent holds the default HTTP User-Agent included in HTTP requests.
//
// Its initial value is what is set below.
//
// Its value can be changed by calling xrpc.SetUserAgent()
var useragent string = "go-xrpc/0.0 (+https://github.com/reiver/go-xrpc)"

// SetUserAgent sets the User-Agent in XRPC requests.
//
// (A [Client] can instead have its own — see [Client].UserAgent.)
func SetUserAgent(value string) {
	useragent = value
}
