package server

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	"github.com/nu7hatch/gouuid"
)

type GameServer struct {
	*gogue.Map
	net.Listener
	Clients []*Client
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

func (gs *GameServer) WaitForClients() {
	for {
		if conn, err := gs.Accept(); err == nil {
			uuid, _ := uuid.NewV4()
			newClient := &Client{
				Player: gs.Game.AddPlayer(uuid.String(), gogue.Position{1, 1, 0}),
				Conn:   conn,
			}
			packet := protocol.Creature{protocol.Position{1, 1, 0}}
			for _, client := range gs.Clients {
				client.Send(packet)
			}

			gs.Clients = append(gs.Clients, newClient)
			newClient.Send(protocol.MapPortion{newClient.Player.MapSight()})
		} else {
			fmt.Println("failed: ", err)
		}
	}
}
