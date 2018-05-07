package caller

import (
	"encoding/json"
	"io"
)

var EmptyParams = []interface{}{}

type APIID uint8

type Caller interface {
	Call(api APIID, method string, args []interface{}, reply interface{}) error
	SetCallback(api APIID, method string, callback func(raw json.RawMessage)) error
}

type CallCloser interface {
	Caller
	io.Closer
}
