package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	"github.com/nu7hatch/gouuid"
)

type Client struct {
	*gogue.Player
	net.Conn
}

type GameServer struct {
	*gogue.Game
	net.Listener
	Clients []*Client
}

func NewGameServer(game *gogue.Game, l net.Listener) (gs *GameServer) {
	return &GameServer{
		Game:     game,
		Listener: l,
	}
}

func Send(c *Client, p protocol.Packet) (err error) {
	enc := gob.NewEncoder(c)
	err = enc.Encode(&p)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return
}

func (gs *GameServer) WaitForClients() {
	gob.Register(protocol.Creature{})
	gob.Register(protocol.MapPortion{})
	for {
		if conn, err := gs.Accept(); err == nil {
			uuid, _ := uuid.NewV4()
			newClient := &Client{
				Player: gs.Game.AddPlayer(uuid.String(), gogue.Position{1, 1, 0}),
				Conn:   conn,
			}
			packet := protocol.Creature{protocol.Position{1, 1, 0}}
			for _, client := range gs.Clients {
				Send(client, packet)
			}

			gs.Clients = append(gs.Clients, newClient)
			Send(newClient, protocol.MapPortion{newClient.Player.MapSight()})
		} else {
			fmt.Println("failed: ", err)
		}
	}
}
