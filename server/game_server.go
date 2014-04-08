package server

import (
	"container/list"
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
)

type GameServer struct {
	*gogue.Map
	net.Listener
	Clients    *list.List
	Broadcast  chan protocol.Packet
	newClients chan *Client
}

func NewGameServer(m *gogue.Map, l net.Listener) (gs *GameServer) {
	return &GameServer{
		Map:        m,
		Listener:   l,
		Clients:    list.New(),
		Broadcast:  make(chan protocol.Packet),
		newClients: make(chan *Client),
	}
}

func (gs *GameServer) Run() {
	go gs.handleClients()

	for {
		if conn, err := gs.Accept(); err == nil {
			client := NewClient(gs.Map, gs.Broadcast)
			na := protocol.NewNetworkAdapter(client.Incoming, client.Outgoing, conn)
			gs.newClients <- client

			go client.Run()
			go na.Listen()
		} else {
			log.Print("failed: ", err)
		}
	}
}

func (gs *GameServer) handleClients() {
	for {
		select {
		case packet := <-gs.Broadcast:
			for e := gs.Clients.Front(); e != nil; e = e.Next() {
				cl := e.Value.(*Client)
				cl.Outgoing <- packet
			}
		case client := <-gs.newClients:
			gs.Clients.PushBack(client)
		}
	}
}
