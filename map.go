package gogue

import "fmt"
import "strings"
import "errors"

type Position struct { X, Y int }

type Player struct { Position }
type Goal struct { Position }

type Tile rune

func (t Tile) String() string {
  return fmt.Sprintf("%c", t)
}

type Map struct {
  Height int
  Width int
  Player
  Goal
  tiles [][]Tile
}

func NewMap(input string) (m Map, err error) {
  var tiles [][]Tile
  var player Player
  var goal Goal
  var playerFound, goalFound bool

  lines := strings.Split(strings.TrimSpace(input), "\n")

  for y, line := range lines {
    tiles = append(tiles, []Tile(strings.TrimSpace(line)))

    for x, tile := range tiles[y] {
      switch tile {
        case Tile('@'):
          player = Player{Position{x, y}}
          playerFound = true
        case Tile('*'):
          goal = Goal{Position{x, y}}
          goalFound = true
      }
    }
  }

  m.Height = len(tiles)
  m.Width = len(tiles[0])
  m.tiles = tiles

  if !playerFound {
    err = errors.New("Map requires a Player(@)")
  } else if !goalFound {
    return m, errors.New("Map requires a Goal(*)")
  } else {
    m.Player = player
    m.Goal = goal
  }

  return
}

  func (m Map) Get(x, y int) Tile {
    return m.tiles[y][x]
  }

