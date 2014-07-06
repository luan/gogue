package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math"
	"net"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	"github.com/onsi/gomega/format"
	"gopkg.in/qml.v0"
)

type tileKey struct {
	gogue.Position
	layer int
}

var creatures = make(map[string]protocol.Creature)
var qmlPlayers = make(map[string]qml.Object)
var mmap gogue.Map
var playerUUID string
var host *string
var tiles = make(map[tileKey]qml.Object)
var lastZ = -99

func inRange(center, pos gogue.Position, r float64) bool {
	if center.Z != pos.Z {
		return false
	}
	diff := center.Diff(pos)
	if math.Abs(float64(diff.X)) > r {
		return false
	}
	if math.Abs(float64(diff.Y)) > r {
		return false
	}
	return true
}

func showMapSight(parent qml.Object) {
	if playerUUID == "" {
		return
	}
	player := creatures[playerUUID]
	z := player.Z
	pos := gogue.Position(player.Position)
	for k, t := range tiles {
		if !inRange(pos, k.Position, 8) {
			t.Set("visible", false)
		}
	}
	for k, p := range qmlPlayers {
		cpos := gogue.Position(creatures[k].Position)
		if !inRange(pos, cpos, 8) {
			p.Set("visible", false)
		}
	}
	lastZ = z
	minX := player.X - 8
	minY := player.Y - 8
	maxX := player.X + 8
	maxY := player.Y + 8
	if minX < 0 {
		minX = 0
	}
	if maxX > mmap.Width {
		maxX = mmap.Width
	}
	if minY < 0 {
		minY = 0
	}
	if maxY >= mmap.Height {
		maxY = mmap.Height
	}
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			tile, _ := mmap.Get(gogue.Position{x, y, z})
			for i, tileID := range tile.Tiles {
				if tileID > 0 {
					key := tileKey{Position: gogue.Position{X: x, Y: y, Z: z}, layer: i}
					if _, ok := tiles[key]; !ok {
						tiles[key] = createTile(tileID, x-minX, y-minY, parent)
					} else {
						tiles[key].Set("x", (x-minX)*tileWidth)
						tiles[key].Set("y", (y-minY)*tileWidth)
						tiles[key].Set("visible", true)
					}
					if tileID >= 1 && tileID < 5 {
						tiles[key].Set("z", i+2)
					}
				}
			}
		}
	}

	for _, cr := range creatures {
		if cr.Z == z {
			x := cr.X - minX
			y := cr.Y - minY
			if _, ok := qmlPlayers[cr.UUID]; !ok {
				qmlPlayers[cr.UUID] = createPlayer(x, y, parent)
			} else {
				qmlPlayers[cr.UUID].Set("x", x*tileWidth)
				qmlPlayers[cr.UUID].Set("y", y*tileWidth)
				qmlPlayers[cr.UUID].Set("visible", true)
			}
		}
	}
}

// func eventLoop(out chan<- protocol.Packet, quit chan bool) {
// 	defer func() {
// 		out <- protocol.Quit{}
// 		close(quit)
// 	}()
// }

type Game struct {
	In   chan protocol.Packet
	Out  chan protocol.Packet
	Quit chan bool
}

type Tile struct {
	Component qml.Object
	Width     int
}

func (g *Game) MoveLeft() {
	g.Out <- protocol.WalkWest
}

func (g *Game) MoveRight() {
	g.Out <- protocol.WalkEast
}

func (g *Game) MoveUp() {
	g.Out <- protocol.WalkNorth
}

func (g *Game) MoveDown() {
	g.Out <- protocol.WalkSouth
}

var tileComponent qml.Object

const tileWidth = 32

func (g *Game) NetLoop(parent qml.Object) {
	conn, err := net.Dial("tcp", *host+":8383")

	if err != nil {
		fmt.Println("Cannot conect to the server.")
		return
	}

	na := protocol.NewNetworkAdapter(g.In, g.Out, g.Quit, conn)
	na.Listen()

	defer na.Close()
	defer func() {
		fmt.Println("Bye bye")
	}()

	for {
		select {
		case p := <-g.In:
			switch t := p.(type) {
			case protocol.Map:
				mmap = p.(protocol.Map).Map
			case protocol.Player:
				pl := p.(protocol.Player)
				playerUUID = pl.UUID
			case protocol.Creature:
				cr := p.(protocol.Creature)
				creatures[cr.UUID] = cr
			case protocol.RemoveCreature:
				cr := p.(protocol.RemoveCreature)
				delete(creatures, cr.UUID)
			case protocol.Quit:
				return
			default:
				log.Print("received unknown packet: ", format.Object(t, 1))
			}

			showMapSight(parent)
		case <-g.Quit:
			return
		}
	}
}

func (g *Game) StartNewGame(parent qml.Object) {
	g.In = make(chan protocol.Packet)
	g.Out = make(chan protocol.Packet)
	g.Quit = make(chan bool)
	go g.NetLoop(parent)
}

func createPlayer(col, row int, parent qml.Object) qml.Object {
	tile := tileComponent.Create(nil)
	tile.Set("parent", parent)

	tile.Set("x", col*tileWidth)
	tile.Set("y", row*tileWidth)
	tile.Set("z", 2)
	tile.Set("width", tileWidth)
	tile.Set("height", tileWidth)
	tile.Set("source", "image://player/1")

	return tile
}

func createTile(tileID, col, row int, parent qml.Object) qml.Object {
	tile := tileComponent.Create(nil)
	tile.Set("parent", parent)

	tile.Set("x", col*tileWidth)
	tile.Set("y", row*tileWidth)
	tile.Set("width", tileWidth)
	tile.Set("height", tileWidth)
	tile.Set("source", "image://sprites/"+strconv.Itoa(tileID))

	return tile
}

func main() {
	host = flag.String("a", "localhost", "Gogue's server address")
	flag.Parse()

	fmt.Println("loading images...")
	images := []image.Image{}
	for i := 1; i <= 7; i++ {
		firstGID := 1
		width := 32
		height := 32
		filename := "../assets/"
		if i >= 6 {
			firstGID = 6
			filename += "stairs.png"
			width = 64
			height = 64
		} else if i >= 5 {
			firstGID = 5
			filename += "ground.png"
		} else {
			width = 64
			height = 64
			filename += "walls.png"
		}
		img, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		start := i - firstGID
		images = append(images, imaging.Crop(img, image.Rect(0, start*height, width, (start+1)*height)))
	}
	playerImage, err := imaging.Open("../assets/player.png")
	if err != nil {
		panic(err)
	}
	fmt.Println("loaded images")

	qml.Init(nil)
	engine := qml.NewEngine()

	engine.AddImageProvider("sprites", func(id string, width, height int) image.Image {
		gid, _ := strconv.Atoi(id)
		return images[gid-1]
	})
	engine.AddImageProvider("player", func(id string, width, height int) image.Image {
		return playerImage
	})

	component, err := engine.LoadFile("gogue.qml")
	if err != nil {
		panic(err)
	}

	game := Game{}
	context := engine.Context()
	context.SetVar("game", &game)

	window := component.CreateWindow(nil)

	tileComponent, err = engine.LoadFile("Tile.qml")
	if err != nil {
		panic(err)
	}

	window.Show()
	window.Wait()
}
