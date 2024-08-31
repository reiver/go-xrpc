# go-xrpc

Package **xrpc** provides an implementation of the **XRPC** protocol used by **BlueSky** and its **AT-Protocol**, for the Go programming language.

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-xrpc

[![GoDoc](https://godoc.org/github.com/reiver/go-xrpc?status.svg)](https://godoc.org/github.com/reiver/go-xrpc)

# Explanation

XRPC is a client-server protocol.
This package implements both the client-side and the server-side of the protocol.

XRPC has 3 requests types:

* `executve` (called `procedure` in the XRPC documentation),
* `query`, and
* `subscribe`.

This package provides functions for making each of these requests:

* `xrpc.Execute()`
* `xrpc.Query()`
* `xrpc.Subscribe()`

(See _examples_ to see how to use each.)

An XRPC URL needs to be passed to each of these functions.
For example:

```golang
const url string = "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"

err := xrpc.Query(dst, url)
```

This package (introduces and) supports 2 types of XRPC URLs:

* `xrpc`, and
* `xrpc-unencrypted`.

For example:

* `xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social`
* `xrpc-unencrypted://localhost/app.bsky.actor.getProfile?actor=reiver.bsky.social`

## Examples

Here is an example of making a 'query' XRPC request:

```golang
import "github.com/reiver/go-xrpc"

// ...

var response map[string]any = map[string]any{}

url := "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"

err := xrpc.Query(&response, url)
```

## XRPC URL Resolution

This package introduces two URL schemes:

* `xrpc`, and
* `xrpc-unencrypted`.

For example:

* `xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social`
* `xrpc-unencrypted://localhost/app.bsky.actor.getProfile?actor=reiver.bsky.social`

Behind-the-scenes the scenes, these XRPC URL schemes are resolved to HTTPS, HTTP, WSS, and WS URLs.

MOST DEVELOPERS WHO ARE JUST MAKING XRPC CLIENT REQUESTS DO NOT HAVE TO WORRY ABOUT THE DETAILS OF THE RESOLUTION.
IN THE SAME WAY THAT MOST DEVELOPERS DO NOT HAVE TO WORRY ABOUT HOW HTTP URLS ARE RESOLVED TO TCP.

How an XRPC URL gets resolved to an HTTPS, HTTP, WSS, or WS URL, depends on the XRPC request type.

Here are some examples:

| XRPC URL                                          | XRPC Request Type | Resolved URL                                 |
|---------------------------------------------------|-------------------|---------------------------------------------:|
| `xrpc://example.com/app.cherry.fooBar`            | `execute`         | `https://example.com/xrpc/app.cherry.fooBar` |
| `xrpc://example.com/app.cherry.fooBar`            | `query`           | `https://example.com/xrpc/app.cherry.fooBar` |
| `xrpc://example.com/app.cherry.fooBar`            | `subscribe`       |   `wss://example.com/xrpc/app.cherry.fooBar` |
| `xrpc-unencrypted://localhost/link.banana.bazQux` | `execute`         |   `http://localhost/xrpc/link.banana.bazQux` |
| `xrpc-unencrypted://localhost/link.banana.bazQux` | `query`           |   `http://localhost/xrpc/link.banana.bazQux` |
| `xrpc-unencrypted://localhost/link.banana.bazQux` | `subscribe`       |     `ws://localhost/xrpc/link.banana.bazQux` |

These 2 XRPC URLs are passed to the `xrpc.Execute()`, `xrpc.Query()`, and `xrpc.Subscribe()` functions.

## Import

To import package **xrpc** use `import` code like the follownig:
```
import "github.com/reiver/go-xrpc"
```

## Installation

To install package **xrpc** do the following:
```
GOPROXY=direct go get https://github.com/reiver/go-xrpc
```

## Author

Package **xrpc** was written by [Charles Iliya Krempeaux](http://reiver.link)
