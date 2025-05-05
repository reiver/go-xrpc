package xrpc

import (
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-xrpcuri"
)

const (
	requestTypeExecute   string = xrpcuri.RequestTypeExecute
	requestTypeQuery     string = xrpcuri.RequestTypeQuery
	requestTypeSubscribe string = xrpcuri.RequestTypeSubscribe
)

func resolve(url string, requestType string) (string, error) {

	resolved, err := xrpcuri.Resolve(url, requestType)
	if nil != err {
		return resolved, erorr.Errorf("xrpc: could not resolve URL %q: %s", url, err)
	}

	return resolved, nil
}
