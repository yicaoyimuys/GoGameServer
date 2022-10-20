package socket

import (
	"GoGameServer/core/libs/sessions"
	"encoding/binary"
	"io"
	"net"
)

func NewFrontCodec(rw net.Conn) sessions.Codec {
	codec := &frontCodec{
		rw:      rw,
		headBuf: make([]byte, 2),
	}
	return codec
}

type frontCodec struct {
	rw      net.Conn
	headBuf []byte
	bodyBuf []byte
}

func (this *frontCodec) Receive() ([]byte, error) {
	//消息长度
	if _, err := io.ReadFull(this.rw, this.headBuf); err != nil {
		return nil, err
	}
	msgLen := binary.BigEndian.Uint16(this.headBuf)

	//消息内容
	if uint16(cap(this.bodyBuf)) < msgLen {
		this.bodyBuf = make([]byte, msgLen, msgLen+128)
	}
	msgBody := this.bodyBuf[:msgLen]
	if _, err := io.ReadFull(this.rw, msgBody); err != nil {
		return nil, err
	}

	return msgBody, nil
}

func (this *frontCodec) Send(msg []byte) error {
	msgLen := uint16(len(msg))
	sendMsg := make([]byte, msgLen+2)
	binary.BigEndian.PutUint16(sendMsg[:2], msgLen)
	copy(sendMsg[2:], msg)

	_, err := this.rw.Write(sendMsg)
	return err
}

func (this *frontCodec) Close() error {
	return this.rw.Close()
}
