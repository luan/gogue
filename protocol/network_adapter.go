package protocol

import (
	"encoding/gob"
	"log"
	"net"
)

type NetworkAdapter struct {
	net.Conn
	*gob.Encoder
	*gob.Decoder
	in   chan<- Packet
	out  <-chan Packet
	quit chan bool
}

func NewNetworkAdapter(in chan<- Packet, out <-chan Packet, quit chan bool, conn net.Conn) *NetworkAdapter {
	gob.Register(Creature{})
	gob.Register(MapPortion{})
	gob.Register(RemoveCreature{})

	gob.Register(Walk{})
	gob.Register(Quit{})

	return &NetworkAdapter{
		Conn:    conn,
		Encoder: gob.NewEncoder(conn),
		Decoder: gob.NewDecoder(conn),
		in:      in,
		out:     out,
		quit:    quit,
	}
}

func (na *NetworkAdapter) Listen() {
	go na.handleOutgoing()
	go na.handleIncoming()
}

func (na *NetworkAdapter) read(p Packet) bool {
	if err := na.Decode(p); err != nil {
		if err.Error() == "EOF" {
			na.Close()
			return false
		} else {
			log.Print("[protocol.NetworkAdapter] decode error:", err)
		}
	}
	return true
}

func (na *NetworkAdapter) handleIncoming() {
	var p Packet
	for na.read(&p) {
		na.in <- p
	}
}

func (na *NetworkAdapter) handleOutgoing() {
	defer na.Close()

	for {
		select {
		case p := <-na.out:
			if err := na.Encode(&p); err != nil {
				log.Print("[protocol.NetworkAdapter] encode error:", err)
			}
		case <-na.quit:
			return
		}
	}
}
