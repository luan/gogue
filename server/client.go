package server

import (
	"log"

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
	go c.listen()
	c.init()
}

func (c *Client) CreaturePacket() protocol.Creature {
	return protocol.Creature{
		UUID:     c.Player.UUID,
		Position: protocol.Position(c.Player.Position),
	}
}

func (c *Client) init() {
	go func() {
		c.Outgoing <- protocol.MapPortion{c.Player.MapSight()}
		c.Outgoing <- c.CreaturePacket()
	}()
	go func() {
		c.Broadcast <- c.CreaturePacket()
	}()
}

func (c *Client) listen() {
	for {
		p := <-c.Incoming
		switch t := p.(type) {
		case protocol.Walk:
			c.processWalk(p.(protocol.Walk))
		default:
			log.Print("received unknown packet: ", t)
		}
	}
}

func (c *Client) processWalk(p protocol.Walk) {
	switch p.Direction {
	case protocol.North:
		c.Player.MoveNorth()
	case protocol.South:
		c.Player.MoveSouth()
	case protocol.East:
		c.Player.MoveEast()
	case protocol.West:
		c.Player.MoveWest()
	}
	c.Broadcast <- c.CreaturePacket()
	c.Outgoing <- c.CreaturePacket()
}
