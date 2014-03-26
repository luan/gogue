package gogue

import (
	"errors"
	"fmt"
)

type Game struct {
	Map
	Player
}

func (g Game) String() string {
	return fmt.Sprintf("Game{\n\t%s\n\t%s\n}",
		indent(fmt.Sprint(g.Player)),
		indent(fmt.Sprint(g.Map)))
}

func (g Game) FollowPath(steps ...rune) (Game, error) {
	if len(steps) == 0 {
		return g, nil
	}

	var err error

	switch steps[0] {
	case 'n':
		g, err = g.MoveNorth()
	case 's':
		g, err = g.MoveSouth()
	case 'e':
		g, err = g.MoveEast()
	case 'w':
		g, err = g.MoveWest()
	}

	if err != nil {
		return g, err
	}

	return g.FollowPath(steps[1:]...)
}

func (g Game) IsOver() bool {
	return g.Player.Position == g.Goal.Position
}

func (g Game) MoveNorth() (Game, error) {
	return g.moveTo(g.Player.North())
}

func (g Game) MoveSouth() (Game, error) {
	return g.moveTo(g.Player.South())
}

func (g Game) MoveEast() (Game, error) {
	return g.moveTo(g.Player.East())
}

func (g Game) MoveWest() (Game, error) {
	return g.moveTo(g.Player.West())
}

func (g Game) moveTo(pos Position) (Game, error) {
	tile := g.Map.Get(pos)

	if !tile.IsWalkable() {
		err := errors.New("cannot move")
		return g, err
	}

	g.Player.Position = pos
	return g, nil
}
