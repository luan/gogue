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
