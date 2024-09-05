package xrpc

import (
	"github.com/gorilla/websocket"
	"github.com/reiver/go-erorr"
)

type internalIterator struct {
	conn *websocket.Conn
	err error
	closed bool
	message []byte
}

var _ Iterator = &internalIterator{}

func (receiver *internalIterator) Close() error {
	if nil == receiver {
		return errNilReceiver
	}

	if receiver.closed {
		return nil
	}

	var conn *websocket.Conn = receiver.conn
	if nil == conn {
		return errNilWebSocketConnection
	}

	return conn.Close()
}

func (receiver *internalIterator) Decode(dst any) error {
	if nil == receiver {
		return errNilReceiver
	}
	if nil == dst {
		return errNilDestination
	}

	if receiver.closed {
		return errClosed
	}

	switch casted := dst.(type) {
	case *string:
		*casted = string(receiver.message)
		return nil
	case *[]byte:
		*casted = append([]byte(nil), receiver.message...)
		return nil
	default:
		return erorr.Errorf("xrpc: cannot decode into something of type %T", dst)
	}
}

func (receiver *internalIterator) Err() error {
	if nil == receiver {
		return errNilReceiver
	}

	return receiver.err
}

func (receiver *internalIterator) Next() bool {
	if nil == receiver {
		return false
	}
	if nil != receiver.err {
		return false
	}
	if receiver.closed {
		return false
	}

	var conn *websocket.Conn = receiver.conn
	if nil == conn {
		receiver.err = errNilWebSocketConnection
		return false
	}

	wsMessageType, wsMessage, err := conn.ReadMessage()
	if nil != err {
		receiver.err = err
		return false
	}

	if websocket.BinaryMessage != wsMessageType {
		receiver.err = erorr.Errorf("xrpc: expected message-type to be %d (binary) but actually got %d", websocket.BinaryMessage, wsMessageType)
		return false
	}

	receiver.message = wsMessage
	return true
}
