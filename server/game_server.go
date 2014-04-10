package server

import (
	"log"
	"net"
	"sync"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
)

type GameServer struct {
	*gogue.Map
	net.Listener
	clients      map[string]*Client
	Broadcast    chan protocol.Packet
	clientsMutex sync.Mutex
}

func NewGameServer(m *gogue.Map, l net.Listener) (gs *GameServer) {
	return &GameServer{
		Map:       m,
		Listener:  l,
		Broadcast: make(chan protocol.Packet),
		clients:   make(map[string]*Client),
	}
}

func (gs *GameServer) Run() {
	go gs.handleClients()

	for {
		if conn, err := gs.Accept(); err == nil {
			cl := NewClient(gs.Map, gs.Broadcast)
			na := protocol.NewNetworkAdapter(cl.Incoming, cl.Outgoing, cl.Quit, conn)
			gs.clientsMutex.Lock()
			gs.clients[cl.UUID] = cl
			gs.clientsMutex.Unlock()

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
			gs.clientsMutex.Lock()
			for _, cl := range gs.clients {
				select {
				case cl.Outgoing <- packet:
				case <-cl.Quit:
				}
			}
			gs.clientsMutex.Unlock()
		}
	}
}

func (gs *GameServer) handleQuit(cl *Client) {
	<-cl.Quit
	gs.clientsMutex.Lock()
	delete(gs.clients, cl.UUID)
	gs.clientsMutex.Unlock()
	gs.Broadcast <- protocol.RemoveCreature{cl.UUID}
}
