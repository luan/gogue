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
	if !g.Map.Get(g.Player.North()).IsWalkable() {
		err := errors.New("cannot move")
		return g, err
	}

	g.Player.Y--
	return g, nil
}

func (g Game) MoveSouth() (Game, error) {
	if !g.Map.Get(g.Player.South()).IsWalkable() {
		err := errors.New("cannot move")
		return g, err
	}

	g.Player.Y++
	return g, nil
}

func (g Game) MoveEast() (Game, error) {
	if !g.Map.Get(g.Player.East()).IsWalkable() {
		err := errors.New("cannot move")
		return g, err
	}

	g.Player.X++
	return g, nil
}

func (g Game) MoveWest() (Game, error) {
	if !g.Map.Get(g.Player.Position.West()).IsWalkable() {
		err := errors.New("cannot move")
		return g, err
	}

	g.Player.X--
	return g, nil
}
