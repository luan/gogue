package gogue

import (
	"errors"
	"fmt"
)

type Player struct {
	UUID string
	Position
	Map *Map
}

func NewPlayer(uuid string, Map *Map, pos Position) *Player {
	return &Player{
		UUID:     uuid,
		Position: pos,
		Map:      Map,
	}
}

func (p *Player) MapSight() (m string) {
	for _, row := range p.Map.Tiles()[p.Z] {
		for _, t := range row {
			m += t.String()
		}

		m += "\n"
	}
	return
}

func (p *Player) MoveNorth() error {
	return p.moveTo(p.North())
}

func (p *Player) MoveSouth() error {
	return p.moveTo(p.South())
}

func (p *Player) MoveEast() error {
	return p.moveTo(p.East())
}

func (p *Player) MoveWest() error {
	return p.moveTo(p.West())
}

func (p *Player) moveTo(pos Position) (err error) {
	if t := p.Map.Get(pos); t.IsWalkable() {
		switch t.ChangeFloor() {
		case "up":
			p.Position = pos.Up()
		case "down":
			p.Position = pos.Down()
		default:
			p.Position = pos
		}
	} else {
		err = errors.New("cannot move")
	}

	return
}

func (p Player) String() string {
	return fmt.Sprintf("Player{%d,%d,%d}", p.X, p.Y, p.Z)
}
