package server

import (
	"log"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	"github.com/nu7hatch/gouuid"
)

type Client struct {
	*gogue.Player
	*gogue.Map
	Broadcast chan<- protocol.Packet
	Outgoing  chan protocol.Packet
	Incoming  chan protocol.Packet
	Quit      chan bool
}

func NewClient(mmap *gogue.Map, broadcast chan<- protocol.Packet) *Client {
	uuid, _ := uuid.NewV4()
	return &Client{
		Player:    gogue.NewPlayer(uuid.String(), mmap, gogue.Position{X: 1, Y: 1, Z: 0}),
		Map:       mmap,
		Broadcast: broadcast,
		Outgoing:  make(chan protocol.Packet),
		Incoming:  make(chan protocol.Packet),
		Quit:      make(chan bool),
	}
}

func (c *Client) Run() {
	c.init()
	go c.listen()
}

func (c *Client) CreaturePacket() protocol.Creature {
	return protocol.Creature{
		UUID:     c.Player.UUID,
		Position: protocol.Position(c.Player.Position),
	}
}

func (c *Client) init() {
	c.Outgoing <- protocol.Map{*c.Map}
	c.Outgoing <- protocol.Player{c.Player.UUID}
	c.Broadcast <- c.CreaturePacket()
}

func (c *Client) listen() {
	defer func() {
		close(c.Quit)
	}()

	for {
		p := <-c.Incoming
		switch t := p.(type) {
		case protocol.Walk:
			c.processWalk(p.(protocol.Walk))
		case protocol.Quit:
			log.Print("quitting: ", c.UUID)
			return
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
}
