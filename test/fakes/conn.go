package fakes

import (
	"encoding/gob"
	"log"
	"net"
	"time"

	"github.com/luan/gogue/protocol"
)

type Conn struct {
	leftBuffer  net.Conn
	rightBuffer net.Conn
}

func NewConn() *Conn {
	leftBuffer, rightBuffer := net.Pipe()
	return &Conn{
		leftBuffer:  leftBuffer,
		rightBuffer: rightBuffer,
	}
}

func (c *Conn) Close() (err error) {
	return
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
	return c.leftBuffer.Read(bytes)
}

func (c *Conn) Write(bytes []byte) (int, error) {
	return c.leftBuffer.Write(bytes)
}

func (c *Conn) Receive() (p protocol.Packet) {
	dec := gob.NewDecoder(c.rightBuffer)
	err := dec.Decode(&p)
	if err != nil {
		log.Print("decode error:", err)
	}
	return
}
