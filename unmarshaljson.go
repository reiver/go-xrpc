package xrpc

import (
	"encoding/json"
	"io"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

func unmarshalJSON(dst any, reader io.Reader) error {
	if nil == reader {
		return errNilReader
	}

	body, err := io.ReadAll(reader)
	if nil != err {
		return erorr.Wrap(ErrResponse, "problem reading the response body", field.Error("error", err))
	}

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
}
