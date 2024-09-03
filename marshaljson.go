package xrpc

import (
	"bytes"
	"io"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-json"
)

func marshalJSON(src any) (io.ReadCloser, error) {
	p, err := json.Marshal(src)
	if nil != err {
		return nil, erorr.Errorf("xrpc: problem json-marshaling value of type %T: %w", src, err)
	}
	if nil == p {
		return nil, erorr.Errorf("xrpc: problem json-marshaling value of type %T: returned bytes is nil", src)
	}

	reader := bytes.NewReader(p)

	var readCloser io.ReadCloser = io.NopCloser(reader)

	return readCloser, nil
}
