package proxys

import (
	"io"
	"net"

	"github.com/funny/link"
	"github.com/funny/binary"
	"time"
)

type clientAddr struct {
	network string
	addr    string
}

func (addr clientAddr) Network() string {
	return addr.network
}

func (addr clientAddr) String() string {
	return addr.addr
}

type DummyConn struct {
	id           uint64
	proxySession *link.Session
	recvChan     chan []byte
	addr         clientAddr
}

func NewDummyConn(id uint64, network string, addr string, proxySession *link.Session) *DummyConn {
	return &DummyConn{
		id:           id,
		proxySession: proxySession,
		recvChan:     make(chan []byte, 1024),
		addr:         clientAddr{network, addr},
	}
}

func (c *DummyConn) LocalAddr() net.Addr {
	return c.proxySession.Conn().LocalAddr()
}

func (c *DummyConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *DummyConn) Read(msg []byte) (int, error) {
	return 0, nil
}

func (c *DummyConn) Write(msg []byte) (int, error) {
	return 0, nil
}

//自定义写消息
func (c *DummyConn) WriteMsg(msg []byte) error {
	sendMsg := make([]byte, 8+len(msg))
	copy(sendMsg[:2], msg[:2])
	binary.PutUint64LE(sendMsg[2:10], c.id)
	copy(sendMsg[10:], msg[2:])
	return c.proxySession.Send(sendMsg)
}

//自定义读消息
func (c *DummyConn) ReadMsg(msg *[]byte) error {
	data, ok := <-c.recvChan
	if !ok {
		return io.EOF
	}

	msgLen := len(data)
	if int64(cap(*msg)) >= int64(msgLen) {
		*msg = (*msg)[0:msgLen]
	} else {
		*msg = make([]byte, msgLen)
	}
	copy(*msg, data)

	return nil
}

//将消息写入Chan
func (c *DummyConn) PutMsg(msg []byte) {
	msgLen := len(msg)-8
	msgID := binary.GetUint16LE(msg[:2])
	msgBody := msg[10:]

	saveMsg := make([]byte, msgLen)
	binary.PutUint16LE(saveMsg[:2], msgID)
	copy(saveMsg[2:], msgBody)

	c.recvChan <- saveMsg
}

func (c *DummyConn) Close() error {
	close(c.recvChan)
	return nil
}

func (c *DummyConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *DummyConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *DummyConn) SetWriteDeadline(t time.Time) error {
	return nil
}
