# go-xrpc

Package **xrpc** provides an implementation of **BlueSky**'s **AT-Protocol**'s XRPC, for the Go programming language.

This package also introduces a `xrpc` and `xrpc-unencrypted` URLs.
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
