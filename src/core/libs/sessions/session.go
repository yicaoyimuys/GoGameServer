package sessions

import (
	"errors"
)

type closeCallback struct {
	Handler interface{}
	Key     interface{}
	Func    func()
	Next    *closeCallback
}

var ErrClosed = errors.New("session closed")
