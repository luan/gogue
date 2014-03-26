package gogue

import (
	"fmt"
)

type Goal struct{ Position }
type Player struct{ Position }

func (p Player) String() string {
	return fmt.Sprintf("Player{%d,%d}", p.X, p.Y)
}

func (g Goal) String() string {
	return fmt.Sprintf("Goal{%d,%d}", g.X, g.Y)
}
