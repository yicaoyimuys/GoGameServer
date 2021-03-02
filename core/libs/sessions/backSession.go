package sessions

import (
	"GoGameServer/core/libs/grpc/ipc"
	"GoGameServer/core/libs/stack"
	"sync"
	"sync/atomic"
)

type BackSession struct {
	id        string
	sessionId uint64
	stream    *ipc.Stream

	closeFlag          int32
	closeChan          chan int
	closeMutex         sync.Mutex
	firstCloseCallback *closeCallback
	lastCloseCallback  *closeCallback

	recvChan  chan []byte
	recvMutex sync.Mutex

	msgHandle func(session *BackSession, msgBody []byte)
	userId    uint64
}

func NewBackSession(id string, sessionId uint64, stream *ipc.Stream) *BackSession {
	session := &BackSession{
		id:        id,
		sessionId: sessionId,
		stream:    stream,
		recvChan:  make(chan []byte, 100),
		closeChan: make(chan int),
	}
	stream.AddSession(session)
	go session.loop()
	return session
}

func (this *BackSession) SessionID() uint64 {
	return this.sessionId
}

func (this *BackSession) ID() string {
	return this.id
}

func (this *BackSession) UserID() uint64 {
	return this.userId
}

func (this *BackSession) SetUserId(userId uint64) {
	this.userId = userId
}

func (this *BackSession) Receive(data []byte) error {
	this.recvMutex.Lock()
	if this.IsClosed() {
		this.recvMutex.Unlock()
		return ErrClosed
	}

	this.recvChan <- data
	this.recvMutex.Unlock()
	return nil
}

func (this *BackSession) Send(data []byte) error {
	if this.IsClosed() {
		return ErrClosed
	}

	return this.stream.Send([]uint64{this.sessionId}, data)
}

func (this *BackSession) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

func (this *BackSession) Close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		close(this.closeChan)

		this.recvMutex.Lock()
		close(this.recvChan)
		for _ = range this.recvChan {
		}
		this.recvMutex.Unlock()

		this.invokeCloseCallbacks()

		this.stream.RemoveSession(this)
		this.stream = nil

		this.msgHandle = nil
	}
}

func (this *BackSession) AddCloseCallback(handler, key interface{}, callback func()) {
	if this.IsClosed() {
		return
	}

	this.closeMutex.Lock()
	defer this.closeMutex.Unlock()

	newItem := &closeCallback{handler, key, callback, nil}

	if this.firstCloseCallback == nil {
		this.firstCloseCallback = newItem
	} else {
		this.lastCloseCallback.Next = newItem
	}
	this.lastCloseCallback = newItem
}

func (this *BackSession) RemoveCloseCallback(handler, key interface{}) {
	if this.IsClosed() {
		return
	}

	this.closeMutex.Lock()
	defer this.closeMutex.Unlock()

	var prev *closeCallback
	for callback := this.firstCloseCallback; callback != nil; prev, callback = callback, callback.Next {
		if callback.Handler == handler && callback.Key == key {
			if this.firstCloseCallback == callback {
				this.firstCloseCallback = callback.Next
			} else {
				prev.Next = callback.Next
			}
			if this.lastCloseCallback == callback {
				this.lastCloseCallback = prev
			}
			return
		}
	}
}

func (this *BackSession) invokeCloseCallbacks() {
	this.closeMutex.Lock()
	defer this.closeMutex.Unlock()

	for callback := this.firstCloseCallback; callback != nil; callback = callback.Next {
		callback.Func()
	}
}

func (this *BackSession) SetMsgHandle(msgHandle func(session *BackSession, msgBody []byte)) {
	this.msgHandle = msgHandle
}

func (this *BackSession) loop() {
	defer stack.TryError()

	for {
		select {
		case msg, ok := <-this.recvChan:
			if ok {
				if this.msgHandle != nil {
					this.msgHandle(this, msg)
				}
			} else {
				return
			}
		case <-this.closeChan:
			return
		}
	}
}
