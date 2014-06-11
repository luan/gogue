package gogue

import "strconv"

type Tile struct {
	Tiles []int
	Properties
}

func (t Tile) PositionModifier() (modifier Position) {
	if changeX, ok := t.Properties["changeX"]; ok {
		modifier.X, _ = strconv.Atoi(changeX)
	}
	if changeY, ok := t.Properties["changeY"]; ok {
		modifier.Y, _ = strconv.Atoi(changeY)
	}
	if changeZ, ok := t.Properties["changeZ"]; ok {
		modifier.Z, _ = strconv.Atoi(changeZ)
	}
	return
}

func (t Tile) IsWalkable() bool {
	return t.Properties["walkable"] != "false"
}
