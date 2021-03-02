package sessions

import (
	"GoGameServer/core/libs/stack"
	"sync"
	"sync/atomic"
	"time"
)

type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}

type FrontSession struct {
	id        uint64
	codec     Codec
	recvMutex sync.Mutex
	sendMutex sync.RWMutex
	recvChan  chan []byte

	closeFlag          int32
	closeChan          chan int
	closeMutex         sync.Mutex
	firstCloseCallback *closeCallback
	lastCloseCallback  *closeCallback

	msgHandle func(session *FrontSession, msgBody []byte)

	pingTime    int64
	ipcServices sync.Map
}

func NewFontSession(id uint64, codec Codec) *FrontSession {
	session := &FrontSession{
		id:        id,
		codec:     codec,
		recvChan:  make(chan []byte, 100),
		closeChan: make(chan int),
		pingTime:  time.Now().Unix(),
	}
	go session.loop()
	return session
}

func (this *FrontSession) ID() uint64 {
	return this.id
}

func (this *FrontSession) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

func (this *FrontSession) GetIpcService(serviceName string) string {
	value, _ := this.ipcServices.Load(serviceName)
	if value != nil {
		return value.(string)
	}
	return ""
}

func (this *FrontSession) PingTime() int64 {
	return this.pingTime
}

func (this *FrontSession) SetIpcService(serviceName string, service string) {
	this.ipcServices.Store(serviceName, service)
}

func (this *FrontSession) UpdatePingTime() {
	this.pingTime = time.Now().Unix()
}

func (this *FrontSession) Close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		close(this.closeChan)

		this.recvMutex.Lock()
		close(this.recvChan)
		for _ = range this.recvChan {
		}
		this.recvMutex.Unlock()

		this.invokeCloseCallbacks()
		this.codec.Close()

		this.msgHandle = nil
	}
}

func (this *FrontSession) Receive() (interface{}, error) {
	msg, err := this.codec.Receive()
	if msg != nil {
		this.recvMutex.Lock()
		if this.IsClosed() {
			this.recvMutex.Unlock()
			return nil, ErrClosed
		}
		this.recvChan <- msg.([]byte)
		this.recvMutex.Unlock()
	}
	return msg, err
}

func (this *FrontSession) Send(msg interface{}) (err error) {
	if this.IsClosed() {
		return ErrClosed
	}

	this.sendMutex.Lock()
	defer this.sendMutex.Unlock()

	return this.codec.Send(msg)
}

func (this *FrontSession) AddCloseCallback(handler, key interface{}, callback func()) {
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

func (this *FrontSession) RemoveCloseCallback(handler, key interface{}) {
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

func (this *FrontSession) invokeCloseCallbacks() {
	this.closeMutex.Lock()
	defer this.closeMutex.Unlock()

	for callback := this.firstCloseCallback; callback != nil; callback = callback.Next {
		callback.Func()
	}
}

func (this *FrontSession) SetMsgHandle(msgHandle func(session *FrontSession, msgBody []byte)) {
	this.msgHandle = msgHandle
}

func (this *FrontSession) loop() {
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
