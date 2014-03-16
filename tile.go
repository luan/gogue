package gogue

import "fmt"
type Tile rune

func (t Tile) String() string {
  return fmt.Sprintf("%c", t)
}

func (t Tile) IsWalkable() bool {
  return t == Tile('.') || t == Tile('*')
}
