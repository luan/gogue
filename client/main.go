package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	termbox "github.com/nsf/termbox-go"
	"github.com/onsi/gomega/format"
)

var creatures = make(map[string]protocol.Creature)
var mmap gogue.Map
var playerUUID string

func showMapSight() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if playerUUID == "" {
		return
	}
	player := creatures[playerUUID]
	z := player.Z
	minX := player.X - 7
	minY := player.Y - 7
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
			t := ' '
			if !tile.IsWalkable() {
				t = '#'
			} else if tile.Tiles[0] == 6 {
				t = '<'
			} else if tile.Tiles[0] == 7 {
				t = '>'
			} else if tile.IsWalkable() {
				t = '.'
			}
			fgAtts := termbox.ColorWhite
			bgAttrs := termbox.ColorDefault

			switch t {
			case '>':
				fgAtts = termbox.ColorBlue
			case '<':
				fgAtts = termbox.ColorCyan
			}

			termbox.SetCell(x-minX+2, y-minY+2, t, fgAtts, bgAttrs)
		}
	}

	for _, cr := range creatures {
		if cr.Z == z {
			fgAtts := termbox.AttrBold + termbox.ColorGreen
			bgAttrs := termbox.ColorDefault
			termbox.SetCell(cr.X-minX+2, cr.Y-minY+2, '@', fgAtts, bgAttrs)
		}
	}
}

func eventLoop(out chan<- protocol.Packet, quit chan bool) {
	defer func() {
		out <- protocol.Quit{}
		close(quit)
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				return
			case termbox.KeyArrowUp:
				out <- protocol.WalkNorth
			case termbox.KeyArrowDown:
				out <- protocol.WalkSouth
			case termbox.KeyArrowLeft:
				out <- protocol.WalkWest
			case termbox.KeyArrowRight:
				out <- protocol.WalkEast
			}
		}
	}
}

func main() {
	host := flag.String("a", "localhost", "Gogue's server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *host+":8383")

	if err != nil {
		fmt.Println("Cannot conect to the server.")
		return
	}

	in := make(chan protocol.Packet)
	out := make(chan protocol.Packet)
	quit := make(chan bool)
	na := protocol.NewNetworkAdapter(in, out, quit, conn)
	na.Listen()

	defer na.Close()

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		termbox.Close()
		fmt.Println("Bye bye")
	}()

	termbox.HideCursor()
	termbox.Flush()

	go eventLoop(out, quit)

	for {
		select {
		case p := <-in:
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

			showMapSight()

			termbox.Flush()
		case <-quit:
			return
		}
	}
}
