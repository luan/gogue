package gogue

import (
  "errors"
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

