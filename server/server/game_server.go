package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/luan/gogue/protocol"
)

type Client interface {
	net.Conn
}

type GameServer struct {
	net.Listener
	Clients []Client
}

func NewGameServer(l net.Listener) (gs *GameServer) {
	return &GameServer{
		Listener: l,
	}
}

func Send(c Client, p protocol.Packet) (err error) {
	enc := gob.NewEncoder(c)
	err = enc.Encode(&p)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return
}

func (gs *GameServer) WaitForClients() {
	gob.Register(protocol.Creature{})
	for {
		if newClient, err := gs.Accept(); err == nil {
			packet := protocol.Creature{protocol.Position{1, 1, 0}}
			for _, client := range gs.Clients {
				Send(client, packet)
			}

			gs.Clients = append(gs.Clients, newClient)
		} else {
			fmt.Println("failed: ", err)
		}
	}
}
