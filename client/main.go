package main

import (
	"fmt"
	"net"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

func showMapSight(mapString string) {
	for y, row := range strings.Split(mapString, "\n") {
		for x, tile := range []byte(row) {
			t := rune(tile)
			fgAtts := termbox.ColorWhite
			bgAttrs := termbox.ColorDefault

			switch t {
			case '>':
				fgAtts = termbox.ColorBlue
			case '<':
				fgAtts = termbox.ColorCyan
			case '*':
				fgAtts = termbox.AttrBold + termbox.ColorYellow
				bgAttrs = termbox.ColorBlack
			case '@':
				fgAtts = termbox.AttrBold + termbox.ColorGreen
			}

			termbox.SetCell(x+5, y+5, t, fgAtts, bgAttrs)
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.HideCursor()
	termbox.Flush()

	conn, err := net.Dial("tcp", "127.0.0.1:8383")

	if err != nil {
		fmt.Println("Cannot conect to the server.")
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error ocurred.")
			return
		}

		if string(buf) == "over" {
			fmt.Println("Congratz, you found the end of the maze!")
			return
		}

		showMapSight(string(buf))
		termbox.Flush()

	ioloop:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyCtrlQ:
					return
				case termbox.KeyArrowUp:
					conn.Write([]byte("mn"))
					break ioloop
				case termbox.KeyArrowDown:
					conn.Write([]byte("ms"))
					break ioloop
				case termbox.KeyArrowLeft:
					conn.Write([]byte("mw"))
					break ioloop
				case termbox.KeyArrowRight:
					conn.Write([]byte("me"))
					break ioloop
				}
			}
		}
	}
}
