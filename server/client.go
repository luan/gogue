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
	return &Client{
		Player:    gogue.NewPlayer(uuid.String(), mmap, gogue.Position{1, 1, 0}),
		Broadcast: broadcast,
		Outgoing:  make(chan protocol.Packet),
		Incoming:  make(chan protocol.Packet),
	}
}

func (c *Client) Run() {
	c.init()
}

func (c *Client) CreaturePacket() protocol.Packet {
	return protocol.Creature{
		protocol.Position(c.Player.Position),
	}
}

func (c *Client) init() {
	go func() {
		c.Outgoing <- protocol.MapPortion{c.Player.MapSight()}
	}()
	go func() {
		c.Broadcast <- c.CreaturePacket()
	}()
}
