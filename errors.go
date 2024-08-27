package xrpc

import (
	"github.com/reiver/go-erorr"
)

const (
	errEmptyHost           = erorr.Error("xrpc: empty host")
	errNilURL              = erorr.Error("xrpc: nil url")
)
