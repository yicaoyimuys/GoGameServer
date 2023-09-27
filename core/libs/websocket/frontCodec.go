package websocket

import (
	"encoding/binary"

	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"

	"github.com/gorilla/websocket"
)

func NewFrontCodec(rw *websocket.Conn) sessions.Codec {
	codec := &frontCodec{
		rw: rw,
	}
	return codec
}

type frontCodec struct {
	rw *websocket.Conn
}

func (this *frontCodec) Receive() ([]byte, error) {
	_, data, err := this.rw.ReadMessage()
	if err != nil {
		return nil, err
	}

	if len(data) < 2 {
		logger.Error("消息长度不够")
		return nil, err
	}

	//消息长度
	msgLen := binary.BigEndian.Uint16(data[:2])
	//消息内容
	msgBody := data[2:]
	//长度检测
	if len(msgBody) != int(msgLen) {
		logger.Error("消息长度不够")
		return nil, err
	}

	return msgBody, nil
}

func (this *frontCodec) Send(msg []byte) error {
	msgLen := uint16(len(msg))
	sendMsg := make([]byte, msgLen+2)
	binary.BigEndian.PutUint16(sendMsg[:2], msgLen)
	copy(sendMsg[2:], msg)

	return this.rw.WriteMessage(websocket.BinaryMessage, sendMsg)
}

func (this *frontCodec) Close() error {
	return this.rw.Close()
}
