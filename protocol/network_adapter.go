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
	in  chan<- Packet
	out <-chan Packet
}

func NewNetworkAdapter(in chan<- Packet, out <-chan Packet, conn net.Conn) *NetworkAdapter {
	gob.Register(Creature{})
	gob.Register(MapPortion{})

	gob.Register(Walk{})

	return &NetworkAdapter{
		Conn:    conn,
		Encoder: gob.NewEncoder(conn),
		Decoder: gob.NewDecoder(conn),
		in:      in,
		out:     out,
	}
}

func (na *NetworkAdapter) Listen() {
	go na.handleOutgoing()
	go na.handleIncoming()
}

func (na *NetworkAdapter) read(p Packet) bool {
	if err := na.Decode(p); err != nil {
		log.Print("decode error:", err)
		na.Close()
		return false
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
	for {
		p := <-na.out
		if err := na.Encode(&p); err != nil {
			log.Print("encode error:", err)
		}
	}
}
