package gogue

import "fmt"

type Tile rune

var WALKABLE_TILES = []Tile{'.', '*', '>', '<'}

func (t Tile) ChangeFloor() (change string) {
	switch t {
	case Tile('>'):
		change = "down"
	case Tile('<'):
		change = "up"
	}
	return
}

func (t Tile) String() string {
	return fmt.Sprintf("%c", t)
}

func (t Tile) IsWalkable() bool {
	return contains(WALKABLE_TILES, t)
}

func contains(s []Tile, t Tile) bool {
	for _, a := range s {
		if a == t {
			return true
		}
	}
	return false
}
