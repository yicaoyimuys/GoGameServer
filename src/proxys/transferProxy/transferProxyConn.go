package transferProxy

import (
	"io"
	"net"
	"github.com/funny/link"
	"github.com/funny/binary"
//	. "tools"
	"time"
)

type clientAddr struct {
	network []byte
	data    []byte
}

func (addr clientAddr) Network() string {
	return string(addr.network)
}

func (addr clientAddr) String() string {
	return string(addr.data)
}

type TransferProxyConn struct {
	id           uint64
	proxySession *link.Session
	recvChan     chan []byte
	addr         clientAddr
}

func NewTransferProxyConn(id uint64, addr clientAddr, proxySession *link.Session) *TransferProxyConn {
	return &TransferProxyConn{
		id:           id,
		proxySession: proxySession,
		recvChan:     make(chan []byte, 1024),
		addr:         addr,
	}
}

func (c *TransferProxyConn) LocalAddr() net.Addr {
	return c.proxySession.Conn().LocalAddr()
}

func (c *TransferProxyConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *TransferProxyConn) Read(msg []byte) (int, error) {
	return 0, nil
}

func (c *TransferProxyConn) ReadOne(msg *[]byte) error {
	data, ok := <-c.recvChan
	if !ok {
		return io.EOF
	}

	msgID := binary.GetUint16LE(data[:2])
	msgBody := data[10:]

	msgLen := len(data)-8
	result := make([]byte, msgLen)
	binary.PutUint16LE(result[:2], msgID)
	copy(result[2:], msgBody)

	if int64(cap(*msg)) >= int64(msgLen) {
		*msg = (*msg)[0:msgLen]
	} else {
		*msg = make([]byte, msgLen)
	}
	copy(*msg, result)

	return nil
}

func (c *TransferProxyConn) Write(msg []byte) (int, error) {
	result := make([]byte, 8+len(msg))
	copy(result[:2], msg[:2])
	binary.PutUint64LE(result[2:10], c.id)
	copy(result[10:], msg[2:])
	c.proxySession.Send(result)
	return 0, nil
}

func (c *TransferProxyConn) Close() error {
	close(c.recvChan)
	return nil
}

func (c *TransferProxyConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *TransferProxyConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *TransferProxyConn) SetWriteDeadline(t time.Time) error {
	return nil
}