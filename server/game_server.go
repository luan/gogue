package server

import (
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
)

type GameServer struct {
	*gogue.Map
	net.Listener
	Clients   map[string]*Client
	Broadcast chan protocol.Packet
}

func NewGameServer(m *gogue.Map, l net.Listener) (gs *GameServer) {
	return &GameServer{
		Map:       m,
		Listener:  l,
		Clients:   make(map[string]*Client),
		Broadcast: make(chan protocol.Packet),
	}
}

func (gs *GameServer) Run() {
	go gs.handleClients()

	for {
		if conn, err := gs.Accept(); err == nil {
			cl := NewClient(gs.Map, gs.Broadcast)
			na := protocol.NewNetworkAdapter(cl.Incoming, cl.Outgoing, cl.Quit, conn)
			gs.Clients[cl.UUID] = cl

			go cl.Run()
			go na.Listen()
			go gs.handleQuit(cl)
		} else {
			log.Print("failed: ", err)
		}
	}
}

func (gs *GameServer) handleClients() {
	for {
		select {
		case packet := <-gs.Broadcast:
			for _, cl := range gs.Clients {
				cl.Outgoing <- packet
			}
		}
	}
}

func (gs *GameServer) handleQuit(cl *Client) {
	<-cl.Quit
	delete(gs.Clients, cl.UUID)
	gs.Broadcast <- protocol.RemoveCreature{cl.UUID}
}
