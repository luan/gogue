package gogue

import "fmt"

type Game struct {
	*Map
	Players map[string]*Player
}

func NewGame(mmap *Map) *Game {
	return &Game{
		Map:     mmap,
		Players: make(map[string]*Player),
	}
}

func (g *Game) AddPlayer(guid string, pos Position) (player *Player) {
	player = NewPlayer(guid, g.Map, pos)
	g.Players[guid] = player
	return
}

func (g Game) String() string {
	return fmt.Sprintf("Game{\n\t%s\n\t%s\n}",
		indent(fmt.Sprint(g.Map)))
}
