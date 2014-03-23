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

func showMapSight(g gogue.Game, light int, wind *ncurses.Window) (s string) {
	for y, row := range g.Tiles() {
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
				wind.Mvaddch(y, x, rune(tile.String()[0]))
			}
		}

		s += "\n"
	}

	return
}

func main() {
	m, _ := gogue.NewMap(`
  ###############################
  #.............................#
  ############........#.....#...#
  #...................#.....#...#
  #........##################...#
  #...................#.........#
  #...##......#########.*.......#
  #...##...#######....#...####..#
  #...##..............#.........#
  #...############....#######...#
  #.............................#
  ###############################
  `)
	game := gogue.Game{
		Map: m,
		Player: gogue.Player{
			Position: gogue.Position{X: 11, Y: 5},
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
	ncurses.Refresh()

	wind := ncurses.NewWindow(game.Height, game.Width, 5, 5)

	for {
		showMapSight(game, 3, wind)
		wind.Refresh()

		if game.IsOver() {
			ncurses.End()
			fmt.Println("Congratz, you found the end of the maze!")
			return
		}

		switch ncurses.Stdscr.Getch() {
		case 'q':
			return
		case 'w':
			game, _ = game.MoveNorth()
		case 's':
			game, _ = game.MoveSouth()
		case 'a':
			game, _ = game.MoveWest()
		case 'd':
			game, _ = game.MoveEast()
		}
	}
}
