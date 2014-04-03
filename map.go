package gogue

import (
	"fmt"
	"strings"
)

type Map struct {
	Height int
	Width  int
	Depth  int
	tiles  [][][]Tile
}

func NewMap(inputs ...string) *Map {
	m := new(Map)
	var tiles [][][]Tile

	for z, input := range inputs {
		lines := strings.Split(strings.TrimSpace(input), "\n")
		tiles = append(tiles, [][]Tile{})

		for _, line := range lines {
			tiles[z] = append(tiles[z], []Tile(strings.TrimSpace(line)))
		}
	}

	m.Depth = len(tiles)
	m.Height = len(tiles[0])
	m.Width = len(tiles[0][0])
	m.tiles = tiles

	return m
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
		indent(m.tilesString()))
}
