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

	switch tile.ChangeFloor() {
	case "up":
		pos = pos.Up()
	case "down":
		pos = pos.Down()
	}
	g.Player.Position = pos
	return g, nil
}
