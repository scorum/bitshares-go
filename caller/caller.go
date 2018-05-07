package caller

import (
	"encoding/json"
	"io"
)

var EmptyParams = []interface{}{}

type Caller interface {
	Call(api uint8, method string, args []interface{}, reply interface{}) error
	SetCallback(api uint8, method string, callback func(raw json.RawMessage)) error
}

type CallCloser interface {
	Caller
	io.Closer
}
