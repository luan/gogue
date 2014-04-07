package server

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	"github.com/nu7hatch/gouuid"
)

type Client struct {
	*gogue.Player
	net.Conn
	Broadcast chan protocol.Packet
}

func NewClient(mmap *gogue.Map, conn net.Conn, broadcast chan protocol.Packet) *Client {
	uuid, _ := uuid.NewV4()
	return &Client{
		Player:    gogue.NewPlayer(uuid.String(), mmap, gogue.Position{1, 1, 0}),
		Conn:      conn,
		Broadcast: broadcast,
	}
}

func (c *Client) Handle() {
	c.Broadcast <- protocol.Creature{protocol.Position(c.Player.Position)}
	c.Send(protocol.MapPortion{c.Player.MapSight()})
}

func (c *Client) Send(p protocol.Packet) (err error) {
	enc := gob.NewEncoder(c)
	err = enc.Encode(&p)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return
}
