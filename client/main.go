package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

func showMapSight(mapString string) {
	mapArray := strings.Split(mapString, "\n")

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
			case '@':
				fgAtts = termbox.AttrBold + termbox.ColorGreen
			case '&':
				fgAtts = termbox.AttrBold + termbox.ColorMagenta
			}

			termbox.SetCell(x+5, y+5, t, fgAtts, bgAttrs)
		}
	}
}

func eventLoop(conn net.Conn) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				conn.Write([]byte("quit"))
			case termbox.KeyArrowUp:
				conn.Write([]byte("mn"))
			case termbox.KeyArrowDown:
				conn.Write([]byte("ms"))
			case termbox.KeyArrowLeft:
				conn.Write([]byte("mw"))
			case termbox.KeyArrowRight:
				conn.Write([]byte("me"))
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

	defer conn.Close()

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

	go eventLoop(conn)

	for {
		buf := make([]byte, 1024)
		bytesRead, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error ocurred.")
			return
		}

		bufString := string(buf[0:bytesRead])

		switch bufString {
		case "quit":
			return
		default:
			showMapSight(bufString)
		}

		termbox.Flush()
	}
}
