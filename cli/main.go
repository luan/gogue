package main

import (
	"fmt"
	"math"

	"github.com/luan/gogue"
	ncurses "github.com/tncardoso/gocurses"
)

func distance(a, b gogue.Position) float64 {
	xDiff := math.Abs(float64(a.X - b.X))
	yDiff := math.Abs(float64(a.Y - b.Y))
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func showMapSight(g gogue.Game, light int, wind *ncurses.Window) {
	tiles := g.Tiles()[g.Player.Z]

	for y, row := range tiles {
		for x, tile := range row {
			if x < 0 || x >= g.Width || y < 0 || y >= g.Height ||
				distance(g.Player.Position, gogue.Position{X: x, Y: y}) > float64(light) {
				wind.Mvaddch(y, x, ' ')
			} else if g.Player.X == x && g.Player.Y == y {
				wind.Attron(ncurses.A_BOLD)
				wind.Attron(ncurses.ColorPair(1))
				wind.Mvaddch(y, x, '@')
				wind.Attroff(ncurses.ColorPair(1))
				wind.Attroff(ncurses.A_BOLD)
			} else {
				t := rune(tile.String()[0])
				if t == '>' {
					wind.Attron(ncurses.ColorPair(2))
				} else if t == '<' {
					wind.Attron(ncurses.ColorPair(3))
				} else if t == '*' {
					wind.Attron(ncurses.ColorPair(4))
				}

				wind.Mvaddch(y, x, t)

				if t == '>' {
					wind.Attroff(ncurses.ColorPair(2))
				} else if t == '<' {
					wind.Attroff(ncurses.ColorPair(3))
				} else if t == '*' {
					wind.Attroff(ncurses.ColorPair(4))
				}
			}
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

	wind := ncurses.NewWindow(game.Height+2, game.Width, 5, 5)

	for {
		showMapSight(game, 7, wind)
		wind.Mvaddstr(game.Height+1, 0, game.Player.String())
		wind.Refresh()

		if game.IsOver() {
			ncurses.End()
			fmt.Println("Congratz, you found the end of the maze!")
			return
		}

		switch ncurses.Stdscr.Getch() {
		case 'q':
			return
		case ncurses.KEY_UP:
			game, _ = game.MoveNorth()
		case ncurses.KEY_DOWN:
			game, _ = game.MoveSouth()
		case ncurses.KEY_LEFT:
			game, _ = game.MoveWest()
		case ncurses.KEY_RIGHT:
			game, _ = game.MoveEast()
		}
	}
}
