package main

import (
	"fmt"
	"math"
	"net"

	"github.com/luan/gogue"
)

var clients uint64

func distance(a, b gogue.Position) float64 {
	xDiff := math.Abs(float64(a.X - b.X))
	yDiff := math.Abs(float64(a.Y - b.Y))
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func sendMapSight(g gogue.Game, light int, conn net.Conn) {
	tiles := g.Tiles()[g.Player.Z]
	mapPacket := ""

	for y, row := range tiles {
		for x, tile := range row {
			if x < 0 || x >= g.Width || y < 0 || y >= g.Height ||
				distance(g.Player.Position, gogue.Position{X: x, Y: y}) > float64(light) {
				mapPacket += " "
			} else if g.Player.X == x && g.Player.Y == y {
				mapPacket += "@"
			} else {
				mapPacket += tile.String()
			}
		}

		mapPacket += "\n"
	}

	conn.Write([]byte(mapPacket))
}

func sendOver(conn net.Conn) {
	conn.Write([]byte("over"))
}

func handleClient(game gogue.Game, conn net.Conn) {
	fmt.Println("Client connected.")
	defer func() {
		fmt.Println("Client left.")
		clients--
		fmt.Printf("Total clients: %d\n", clients)
		conn.Close()
	}()

	for {
		sendMapSight(game, 7, conn)

		if game.IsOver() {
			sendOver(conn)
			return
		}

		buf := make([]byte, 4)
		bytesRead, err := conn.Read(buf)
		if err != nil {
			return
		}

		bufString := string(buf[0:bytesRead])
		fmt.Println("Client(", conn, "): ", bufString)

		switch bufString {
		case "quit":
			conn.Write([]byte("quit"))
		case "mn":
			game, _ = game.MoveNorth()
		case "ms":
			game, _ = game.MoveSouth()
		case "mw":
			game, _ = game.MoveWest()
		case "me":
			game, _ = game.MoveEast()
		}
	}
}

func main() {
	m, _ := gogue.NewMap(`
	###############################
	#.............................#
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
	#*#.###..#....................#
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
	game := gogue.Game{
		Map: m,
		Player: gogue.Player{
			Position: gogue.Position{X: 11, Y: 5, Z: 0},
		},
	}

	ln, err := net.Listen("tcp", ":8383")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening on port 8383")

	for {
		conn, err := ln.Accept()
		clients++
		fmt.Printf("Total clients: %d\n", clients)
		if err != nil {
			fmt.Println(err)
		}

		go handleClient(game, conn)
	}
}
