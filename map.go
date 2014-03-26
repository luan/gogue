package gogue

import (
	"errors"
	"fmt"
	"strings"
)

type Map struct {
	Height int
	Width  int
	Depth  int
	Goal
	tiles [][][]Tile
}

func NewMap(inputs ...string) (m Map, err error) {
	var tiles [][][]Tile
	var goal Goal
	var goalFound bool

	for z, input := range inputs {
		lines := strings.Split(strings.TrimSpace(input), "\n")
		tiles = append(tiles, [][]Tile{})

		for y, line := range lines {
			tiles[z] = append(tiles[z], []Tile(strings.TrimSpace(line)))

			for x, tile := range tiles[z][y] {
				switch tile {
				case Tile('*'):
					goal = Goal{Position{X: x, Y: y, Z: z}}
					goalFound = true
				}
			}
		}
	}

	m.Depth = len(tiles)
	m.Height = len(tiles[0])
	m.Width = len(tiles[0][0])
	m.tiles = tiles

	if !goalFound {
		err = errors.New("map requires a Goal(*)")
	} else {
		m.Goal = goal
	}

	return
}

func (m Map) Tiles() [][][]Tile {
	return m.tiles
}

func (m Map) Get(pos Position) Tile {
	return m.tiles[pos.Z][pos.Y][pos.X]
}

func (m Map) tilesString() (s string) {
	for _, row := range m.tiles {
		for _, tile := range row {
			s += fmt.Sprint(tile)
		}
		s += "\n"
	}

	return
}

func (m Map) String() string {
	return fmt.Sprintf("Map{\n\tDimension{%d,%d}\n\t%s\n\t%s\n}",
		m.Height,
		m.Width,
		fmt.Sprint(m.Goal),
		indent(m.tilesString()))
}
