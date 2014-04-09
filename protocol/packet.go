package protocol

type Packet interface{}

type Creature struct {
	UUID string
	Position
}

type Position struct {
	X, Y, Z int
}

type MapPortion struct {
	Data string
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
