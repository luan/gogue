package gogue

import "errors"

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
	var t Tile
	if t, err = p.Map.Get(pos); err == nil && t.IsWalkable() {
		p.Position = pos.Add(t.PositionModifier())
	} else {
		err = errors.New("cannot move")
	}
	return
}
