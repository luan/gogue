package protocol

import "github.com/luan/gogue"

type Packet interface{}

type Player struct {
	UUID string
}

type Creature struct {
	UUID string
	Position
}

type RemoveCreature struct {
	UUID string
}

type Position struct {
	X, Y, Z int
}

type Map struct {
	Map gogue.Map
}

const (
	North = "n"
	South = "s"
	East  = "e"
	West  = "w"
)

type Walk struct {
	Direction string
}

var (
	WalkNorth = Walk{Direction: North}
	WalkSouth = Walk{Direction: South}
	WalkWest  = Walk{Direction: West}
	WalkEast  = Walk{Direction: East}
)

type Quit struct{}
