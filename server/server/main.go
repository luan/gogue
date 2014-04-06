package main

import (
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/luan/gogue"
)

var totalClients int
var clients map[string]*Client

type Client struct {
	Id string
	*gogue.Game
	*gogue.Player
	net.Conn
}

func distance(a, b gogue.Position) float64 {
	xDiff := math.Abs(float64(a.X - b.X))
	yDiff := math.Abs(float64(a.Y - b.Y))
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func playerAt(game *gogue.Game, pos gogue.Position) bool {
	for _, player := range game.Players {
		if player.X == pos.X && player.Y == pos.Y && player.Z == pos.Z {
			return true
		}
	}

	return false
}

func sendMapSight(client *Client, light int) {
	tiles := client.Game.Tiles()[client.Player.Z]
	mapPacket := ""

	for y, row := range tiles {
		for x, tile := range row {
			if x < 0 || x >= client.Game.Width || y < 0 || y >= client.Game.Height ||
				distance(client.Player.Position, gogue.Position{X: x, Y: y}) > float64(light) {
				mapPacket += " "
			} else if client.Player.X == x && client.Player.Y == y {
				mapPacket += "@"
			} else if playerAt(client.Game, gogue.Position{x, y, client.Player.Z}) {
				mapPacket += "&"
			} else {
				mapPacket += tile.String()
			}
		}

		mapPacket += "\n"
	}

	client.Conn.Write([]byte(mapPacket))
}

func handleClient(client *Client, mapUpdate chan<- bool) {
	fmt.Println("Client connected.")
	defer func() {
		client.Conn.Close()
		delete(client.Game.Players, client.Player.Guid)
		delete(clients, client.Id)

		fmt.Printf("Client[%s] - Left\n", client.Player.Guid)
		fmt.Printf("Total connecetd clients: %d\n", len(clients))
		mapUpdate <- true
	}()

	for {
		buf := make([]byte, 4)
		bytesRead, err := client.Conn.Read(buf)
		if err != nil {
			return
		}

		bufString := string(buf[0:bytesRead])
		fmt.Printf("Client[%s]: `%s`\n", client.Player.Guid, bufString)

		switch bufString {
		case "quit":
			client.Conn.Write([]byte("quit"))
			return
		case "mn":
			client.Player.MoveNorth()
		case "ms":
			client.Player.MoveSouth()
		case "mw":
			client.Player.MoveWest()
		case "me":
			client.Player.MoveEast()
		}
		mapUpdate <- true
	}
}

func handleMapUpdates(clients map[string]*Client, mapUpdate <-chan bool) {
	for {
		<-mapUpdate
		for _, client := range clients {
			sendMapSight(client, 7)
		}
	}
}

func main() {
	m := gogue.NewMap(`
	###############################
	#>............................#
	############........#.....#...#
	#...................#.....#...#
	#........##################...#
	#...................#.........#
	#...##......#########.>.......#
	#...##...#######....#...####..#
	#...##..............#.........#
	#...############....#######...#
	#.............................#
	###############################
	`, `
	###############################
	#<#.###..#....................#
	#.#......#.##################.#
	#.#.#.####.########....######.#
	#.#.#.............#...........#
	#.#.#########################.#
	#.#...................<.....#.#
	#.#.....###########.........#.#
	#.###########################.#
	#..........#...#..............#
	#........#...#...#............#
	###############################
	`)
	game := gogue.NewGame(m)

	ln, err := net.Listen("tcp", ":8383")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening on port 8383")

	mapUpdate := make(chan bool)
	clients = make(map[string]*Client)

	go handleMapUpdates(clients, mapUpdate)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}

		totalClients++
		clientID := strconv.Itoa(totalClients)
		player := game.AddPlayer(clientID, gogue.Position{X: 11, Y: 5, Z: 0})
		client := &Client{clientID, game, player, conn}
		clients[clientID] = client
		fmt.Printf("Total connected clients: %d\n", len(clients))

		go handleClient(client, mapUpdate)
		mapUpdate <- true
	}
}
