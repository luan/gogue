package fakes

import (
	"net"
	"sync"
	"time"
)

type Conn struct {
	buffer      net.Conn
	closed      bool
	closedMutex sync.Mutex
}

func NewConn(buffer net.Conn) *Conn {
	return &Conn{
		buffer: buffer,
	}
}

func (c *Conn) Close() (err error) {
	c.closedMutex.Lock()
	c.closed = true
	c.closedMutex.Unlock()
	return c.buffer.Close()
}

func (c *Conn) Closed() bool {
	c.closedMutex.Lock()
	defer c.closedMutex.Unlock()
	return c.closed
}

func (c *Conn) LocalAddr() (addr net.Addr) {
	return
}

func (c *Conn) RemoteAddr() (addr net.Addr) {
	return
}

func (c *Conn) SetDeadline(t time.Time) (err error) {
	return
}

func (c *Conn) SetReadDeadline(t time.Time) (err error) {
	return
}

func (c *Conn) SetWriteDeadline(t time.Time) (err error) {
	return
}

func (c *Conn) Read(bytes []byte) (int, error) {
	return c.buffer.Read(bytes)
}

func (c *Conn) Write(bytes []byte) (int, error) {
	return c.buffer.Write(bytes)
}
