package gogue

type Position struct{ X, Y int }

func (p Position) North() Position {
	p.Y--
	return p
}

func (p Position) South() Position {
	p.Y++
	return p
}

func (p Position) East() Position {
	p.X++
	return p
}

func (p Position) West() Position {
	p.X--
	return p
}
