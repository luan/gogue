package main

import (
	"fmt"
	"net"
	"strings"

	ncurses "github.com/tncardoso/gocurses"
)

func showMapSight(mapString string, wind *ncurses.Window) {
	for y, row := range strings.Split(mapString, "\n") {
		for x, tile := range []byte(row) {
			t := rune(tile)
			if t == '>' {
				wind.Attron(ncurses.ColorPair(2))
			} else if t == '<' {
				wind.Attron(ncurses.ColorPair(3))
			} else if t == '*' {
				wind.Attron(ncurses.ColorPair(4))
			} else if t == '@' {
				wind.Attron(ncurses.A_BOLD)
				wind.Attron(ncurses.ColorPair(1))
			}

			wind.Mvaddch(y, x, t)

			if t == '>' {
				wind.Attroff(ncurses.ColorPair(2))
			} else if t == '<' {
				wind.Attroff(ncurses.ColorPair(3))
			} else if t == '*' {
				wind.Attroff(ncurses.ColorPair(4))
			} else if t == '@' {
				wind.Attroff(ncurses.A_BOLD)
				wind.Attroff(ncurses.ColorPair(1))
			}
		}
	}
}

func main() {
	ncurses.Initscr()
	defer ncurses.End()

	ncurses.StartColor()
	ncurses.Cbreak()
	ncurses.Noecho()
	ncurses.Stdscr.Keypad(true)

	ncurses.CursSet(0)
	ncurses.InitPair(1, ncurses.COLOR_GREEN, ncurses.COLOR_BLACK)
	ncurses.InitPair(2, ncurses.COLOR_BLUE, ncurses.COLOR_BLACK)
	ncurses.InitPair(3, ncurses.COLOR_CYAN, ncurses.COLOR_BLACK)
	ncurses.InitPair(4, ncurses.COLOR_RED, ncurses.COLOR_BLACK)
	ncurses.Refresh()

	wind := ncurses.NewWindow(12, 31, 5, 5)
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
			ncurses.End()
			fmt.Println("Congratz, you found the end of the maze!")
			return
		}

		showMapSight(string(buf), wind)
		// wind.Mvaddstr(game.Height+1, 0, game.Player.String())
		wind.Refresh()

		switch ncurses.Stdscr.Getch() {
		case 'q':
			return
		case ncurses.KEY_UP:
			conn.Write([]byte("mn"))
		case ncurses.KEY_DOWN:
			conn.Write([]byte("ms"))
		case ncurses.KEY_LEFT:
			conn.Write([]byte("mw"))
		case ncurses.KEY_RIGHT:
			conn.Write([]byte("me"))
		}
	}
}
