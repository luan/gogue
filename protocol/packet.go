package protocol

type Packet interface{}

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

type MapPortion struct {
	Data string
	Z    int
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
