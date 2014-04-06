package server

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
)

type Client struct {
	*gogue.Player
	net.Conn
}

func (c *Client) Send(p protocol.Packet) (err error) {
	enc := gob.NewEncoder(c)
	err = enc.Encode(&p)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return
}
