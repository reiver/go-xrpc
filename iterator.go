package xrpc

// Iterator represents an iterator.
type Iterator interface {
	Close() error
	Decode(any) error
	Err() error
	Next() bool
}
