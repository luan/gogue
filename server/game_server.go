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
			log.Print("new client connected: ", cl.UUID)
			na := protocol.NewNetworkAdapter(cl.Incoming, cl.Outgoing, cl.Quit, conn)
			gs.addClient(cl)

			go cl.Run()
			go na.Listen()
			go gs.handleQuit(cl)

			go func() {
				for _, ocl := range gs.Clients() {
					if cl.Player.Z == ocl.Player.Z {
						cl.Outgoing <- ocl.CreaturePacket()
					}
				}
			}()
		} else {
			log.Print("failed: ", err)
		}
	}
}

func (gs *GameServer) Clients() (clients map[string]*Client) {
	gs.clientsMutex.Lock()
	defer gs.clientsMutex.Unlock()
	clients = make(map[string]*Client)
	for k, v := range gs.clients {
		clients[k] = v
	}
	return
}

func (gs *GameServer) addClient(cl *Client) {
	gs.clientsMutex.Lock()
	defer gs.clientsMutex.Unlock()
	gs.clients[cl.UUID] = cl
}

func (gs *GameServer) deleteClient(UUID string) {
	gs.clientsMutex.Lock()
	defer gs.clientsMutex.Unlock()
	delete(gs.clients, UUID)
}

func (gs *GameServer) handleClients() {
	for {
		select {
		case packet := <-gs.Broadcast:
			for _, cl := range gs.Clients() {
				var scl *Client
				if p, ok := packet.(protocol.RemoveCreature); ok {
					scl = gs.Clients()[p.UUID]
				}
				if p, ok := packet.(protocol.Creature); ok {
					scl = gs.Clients()[p.UUID]
				}

				if scl != nil && cl.Player.Z != scl.Player.Z {
					continue
				}
				log.Print("sending broadcast to: ", cl.UUID)
				select {
				case cl.Outgoing <- packet:
				case <-cl.Quit:
				}
			}
		}
	}
}

func (gs *GameServer) handleQuit(cl *Client) {
	<-cl.Quit
	log.Print("client left: ", cl.UUID)
	gs.deleteClient(cl.UUID)
	gs.Broadcast <- protocol.RemoveCreature{UUID: cl.UUID}
}
