package sessions

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
)

func NewFrontCodec(rw *websocket.Conn) Codec {
	codec := &frontByteCodec{
		rw: rw,
	}
	return codec
}

type frontByteCodec struct {
	rw *websocket.Conn
}

func (this *frontByteCodec) Receive() (interface{}, error) {
	_, data, err := this.rw.ReadMessage()
	if err != nil {
		return nil, err
	}

	//消息长度
	//	msgLen := binary.BigEndian.Uint16(data[:2])

	//消息内容
	msgBody := data[2:]

	return msgBody, nil
}

func (this *frontByteCodec) Send(msg1 interface{}) error {
	msg := msg1.([]byte)

	msgLen := uint16(len(msg))
	sendMsg := make([]byte, msgLen+2)
	binary.BigEndian.PutUint16(sendMsg[:2], msgLen)
	copy(sendMsg[2:], msg)

	return this.rw.WriteMessage(websocket.BinaryMessage, sendMsg)
}

func (this *frontByteCodec) Close() error {
	return this.rw.Close()
}
