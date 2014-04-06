package fakes

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"time"

	"github.com/luan/gogue/protocol"
)

type Client struct {
	bytes.Buffer
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(l *Listener) {
	l.ch <- c
	return
}

func (c *Client) Close() (err error) {
	return
}

func (c *Client) LocalAddr() (addr net.Addr) {
	return
}

func (c *Client) RemoteAddr() (addr net.Addr) {
	return
}

func (c *Client) SetDeadline(t time.Time) (err error) {
	return
}

func (c *Client) SetReadDeadline(t time.Time) (err error) {
	return
}

func (c *Client) SetWriteDeadline(t time.Time) (err error) {
	return
}

func (c *Client) Receive() (p protocol.Packet, err error) {
	if c.Len() > 0 {
		dec := gob.NewDecoder(&c.Buffer)
		err = dec.Decode(&p)
		if err != nil {
			log.Fatal("decode error:", err)
		}
	} else {
		err = errors.New("no data to be read")
	}
	return
}
