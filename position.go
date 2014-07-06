package gogue

type Position struct{ X, Y, Z int }

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

func (p Position) Up() Position {
	p.Z++
	return p
}

func (p Position) Down() Position {
	p.Z--
	return p
}

func (p Position) Diff(q Position) Position {
	return Position{
		X: p.X - q.X,
		Y: p.Y - q.Y,
		Z: p.Z - q.Z,
	}
}

func (p Position) Add(q Position) Position {
	return Position{
		X: p.X + q.X,
		Y: p.Y + q.Y,
		Z: p.Z + q.Z,
	}
}
