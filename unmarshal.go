package xrpc

import (
	"bytes"
	"encoding/json"
	"io"
)

func unmarshal(dst any, reader io.Reader) error {
	if nil == reader {
		return errNilReader
	}

	var bodyBuffer bytes.Buffer
	io.Copy(&bodyBuffer, reader)

	var body []byte = bodyBuffer.Bytes()

	switch casted := dst.(type) {
	case json.Unmarshaler:
		return json.Unmarshal(body, casted)
	case *string:
		*casted = string(body)
		return nil
	case *[]byte:
		*casted = body
		return nil
	case *[]rune:
		*casted = []rune(string(body))
		return nil
	default:
		return json.Unmarshal(body, dst)
	}

	return nil
}
