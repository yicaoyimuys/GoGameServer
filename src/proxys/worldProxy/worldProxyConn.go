package worldProxy

import (
	"bytes"
	"io"
	"net"

	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	//	. "tools"
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

type WorldProxyConn struct {
	id           uint64
	proxySession *link.Session
	recvChan     chan packet.RAW
	addr         clientAddr
}

func NewWorldProxyConn(id uint64, addr clientAddr, proxySession *link.Session) *WorldProxyConn {
	return &WorldProxyConn{
		id:           id,
		proxySession: proxySession,
		recvChan:     make(chan packet.RAW, 1024),
		addr:         addr,
	}
}

func (c *WorldProxyConn) Config() link.SessionConfig {
	return link.SessionConfig{
		1024,
	}
}

func (c *WorldProxyConn) LocalAddr() net.Addr {
	return c.proxySession.Conn().LocalAddr()
}

func (c *WorldProxyConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *WorldProxyConn) Receive(msg interface{}) error {
	data, ok := <-c.recvChan
	if !ok {
		return io.EOF
	}

	msgID := binary.GetUint16LE(data[:2])
	msgBody := data[10:]

	result := make([]byte, len(data)-8)
	binary.PutUint16LE(result[:2], msgID)
	copy(result[2:], msgBody)

	if fast, ok := msg.(packet.FastInMessage); ok {
		return fast.Unmarshal(
			&io.LimitedReader{bytes.NewReader(result), int64(len(result))},
		)
	}

	msg.(packet.InMessage).Unmarshal(result)
	return nil
}

func (c *WorldProxyConn) Send(data interface{}) error {
	msg := data.(packet.RAW)
	result := make([]byte, 8+len(msg))
	copy(result[:2], msg[:2])
	binary.PutUint64LE(result[2:10], c.id)
	copy(result[10:], msg[2:])
	c.proxySession.Send(packet.RAW(result))
	return nil
}

func (c *WorldProxyConn) Close() error {
	close(c.recvChan)
	return nil
}
