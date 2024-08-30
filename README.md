# go-xrpc

Package **xrpc** provides an implementation of the **XRPC** protocol used by **BlueSky** and its **AT-Protocol**, for the Go programming language.

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

An XRPC 'url' needs to be passed to each of these functions.
For example:

```golang
const url string = "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"

err := xrpc.Query(dst, url)
```

This (introduces and) supports 2 types of XRPC URLs:

* `xrpc`, and
* `xrpc-unencrypted`.

For example:

* `xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social`
* `xrpc-unencrypted://localhost/app.bsky.actor.getProfile?actor=reiver.bsky.social`

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-xrpc

[![GoDoc](https://godoc.org/github.com/reiver/go-xrpc?status.svg)](https://godoc.org/github.com/reiver/go-xrpc)

## Example

Here is an example of making a 'query' XRPC request:

```golang
import "github.com/reiver/go-xrpc"

// ...

var response map[string]any = map[string]any{}

url := "xrpc://public.api.bsky.app/app.bsky.actor.getProfile?actor=reiver.bsky.social"

err := xrpc.Query(&response, url)
```

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
