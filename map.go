package gogue

import (
  "errors"
  "fmt"
  "strings"
)

type Map struct {
  Height int
  Width int
  Goal
  tiles [][]Tile
}

func NewMap(input string) (m Map, err error) {
  var tiles [][]Tile
  var goal Goal
  var goalFound bool

  lines := strings.Split(strings.TrimSpace(input), "\n")

  for y, line := range lines {
    tiles = append(tiles, []Tile(strings.TrimSpace(line)))

    for x, tile := range tiles[y] {
      switch tile {
      case Tile('*'):
        goal = Goal{Position{x, y}}
        goalFound = true
      }
    }
  }

  m.Height = len(tiles)
  m.Width = len(tiles[0])
  m.tiles = tiles

  if !goalFound {
    err = errors.New("Map requires a Goal(*)")
  } else {
    m.Goal = goal
  }

  return
}

func (m Map) Get(x, y int) Tile {
  return m.tiles[y][x]
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
