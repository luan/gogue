package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/luan/gogue/protocol"
	termbox "github.com/nsf/termbox-go"
)

var creatures = make(map[string]protocol.Creature)
var mapPortion protocol.MapPortion

func showMapSight() {
	mapArray := strings.Split(mapPortion.Data, "\n")

	for y, row := range mapArray {
		for x, tile := range []byte(row) {
			t := rune(tile)
			fgAtts := termbox.ColorWhite
			bgAttrs := termbox.ColorDefault

			switch t {
			case '>':
				fgAtts = termbox.ColorBlue
			case '<':
				fgAtts = termbox.ColorCyan
			}

			termbox.SetCell(x+5, y+5, t, fgAtts, bgAttrs)
		}
	}

	for _, cr := range creatures {
		if cr.Z == mapPortion.Z {
			fgAtts := termbox.AttrBold + termbox.ColorGreen
			bgAttrs := termbox.ColorDefault
			termbox.SetCell(cr.X+5, cr.Y+5, '@', fgAtts, bgAttrs)
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
			case protocol.MapPortion:
				mapPortion = p.(protocol.MapPortion)
			case protocol.Creature:
				cr := p.(protocol.Creature)
				creatures[cr.UUID] = cr
			case protocol.RemoveCreature:
				cr := p.(protocol.RemoveCreature)
				delete(creatures, cr.UUID)
			default:
				log.Print("received unknown packet: ", t)
			}

			showMapSight()

			termbox.Flush()
		case <-quit:
			return
		}
	}
}
