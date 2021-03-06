package main

import (
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/server"
)

func main() {
	m := gogue.NewMap("../../assets/map-tiled.json")

	ln, err := net.Listen("tcp", ":8383")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Listening on port 8383")
	server.NewGameServer(m, ln).Run()
}
