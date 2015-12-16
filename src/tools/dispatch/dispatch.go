package dispatch

import (
	"github.com/funny/link"
)

type ReceiveMsg struct {
	Session *link.Session
	Msg     []byte
}

type ReceiveMsgChan chan ReceiveMsg

type DispatchInterface interface {
	Process(session *link.Session, msg []byte)
}

//同步Dispatch
type Dispatch struct {
	handle HandleInterface
}

func NewDispatch(h HandleInterface) Dispatch {
	return Dispatch{
		handle: h,
	}
}

func (this Dispatch) Process(session *link.Session, msg []byte) {
	this.handle.DealMsg(session, msg)
}

//异步Dispatch
type DispatchAsync struct {
	handle HandleInterface
	receiveMsgs []ReceiveMsgChan
	receiveMsgsIndex int
	receiveMsgsLen int
}

func NewDispatchAsync(receiveMsgChans []ReceiveMsgChan, h HandleInterface) DispatchAsync {
	dispatch := DispatchAsync{
		handle: h,
		receiveMsgs: receiveMsgChans,
		receiveMsgsIndex: 0,
		receiveMsgsLen: len(receiveMsgChans),
	}
	dispatch.dealReceiveMsgs()
	return dispatch
}

func (this DispatchAsync) Process(session *link.Session, msg []byte) {
	this.receiveMsgs[this.receiveMsgsIndex] <- ReceiveMsg{
		Session: session,
		Msg: msg,
	}

	this.receiveMsgsIndex += 1
	if this.receiveMsgsIndex == len(this.receiveMsgs) {
		this.receiveMsgsIndex = 0
	}
}

func (this DispatchAsync) dealReceiveMsgs() {
	for i := 0; i < this.receiveMsgsLen; i++ {
		receiveMsgChan := this.receiveMsgs[i]
		go func(msgChan ReceiveMsgChan) {
			for {
				data, ok := <-msgChan
				if !ok {
					return
				}
				this.handle.DealMsg(data.Session, data.Msg)
			}
		}(receiveMsgChan)
	}
}