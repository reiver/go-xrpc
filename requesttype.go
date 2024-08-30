package xrpc

const (
	RequstTypeExecute    string = "execute" // an alternative name for "procedure" so that all the request-types can be a verb.
	RequestTypeProcedure string = "procedure"
	RequestTypeQuery     string = "query"
	RequestTypeSubscribe string = "subscribe"
)
