package server

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
)

type GameServer struct {
	*gogue.Map
	net.Listener
	Clients   *list.List
	Broadcast chan protocol.Packet
}

func NewGameServer(m *gogue.Map, l net.Listener) (gs *GameServer) {
	gob.Register(protocol.Creature{})
	gob.Register(protocol.MapPortion{})

	return &GameServer{
		Map:       m,
		Listener:  l,
		Clients:   list.New(),
		Broadcast: make(chan protocol.Packet),
	}
}

func (gs *GameServer) Run() {
	go gs.ListenBroadcast()

	for {
		if conn, err := gs.Accept(); err == nil {
			client := NewClient(gs.Map, conn, gs.Broadcast)
			gs.Clients.PushBack(client)

			client.Handle()
		} else {
			fmt.Println("failed: ", err)
		}
	}
}

func (gs *GameServer) ListenBroadcast() {
	for {
		select {
		case packet := <-gs.Broadcast:
			for e := gs.Clients.Front(); e != nil; e = e.Next() {
				client := e.Value.(*Client)
				client.Send(packet)
			}
		}
	}
}
