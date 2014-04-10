package fakes

import (
	"encoding/gob"
	"net"
	"time"

	"github.com/luan/gogue/protocol"
)

type Conn struct {
	serverBuffer net.Conn
	clientBuffer net.Conn
}

func NewConn() *Conn {
	serverBuffer, clientBuffer := net.Pipe()
	return &Conn{
		serverBuffer: serverBuffer,
		clientBuffer: clientBuffer,
	}
}

func (c *Conn) Close() (err error) {
	return c.clientBuffer.Close()
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
	return c.serverBuffer.Read(bytes)
}

func (c *Conn) Write(bytes []byte) (int, error) {
	return c.serverBuffer.Write(bytes)
}

func (c *Conn) Send(p protocol.Packet) {
	enc := gob.NewEncoder(c.clientBuffer)
	enc.Encode(&p)
}

func (c *Conn) Receive() (p protocol.Packet) {
	dec := gob.NewDecoder(c.clientBuffer)
	dec.Decode(&p)
	return
}
