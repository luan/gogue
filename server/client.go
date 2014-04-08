package server

import (
	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	"github.com/nu7hatch/gouuid"
)

type Client struct {
	*gogue.Player
	Broadcast chan<- protocol.Packet
	Outgoing  chan protocol.Packet
	Incoming  chan protocol.Packet
}

func NewClient(mmap *gogue.Map, broadcast chan<- protocol.Packet) *Client {
	uuid, _ := uuid.NewV4()
	in, out := make(chan protocol.Packet), make(chan protocol.Packet)
	return &Client{
		Player:    gogue.NewPlayer(uuid.String(), mmap, gogue.Position{1, 1, 0}),
		Broadcast: broadcast,
		Outgoing:  out,
		Incoming:  in,
	}
}

func (c *Client) Run() {
	c.init()
}

func (c *Client) init() {
	for i := 0; i < 2; i++ {
		select {
		case c.Outgoing <- protocol.MapPortion{c.Player.MapSight()}:
		case c.Broadcast <- protocol.Creature{protocol.Position(c.Player.Position)}:
		}
	}
}
