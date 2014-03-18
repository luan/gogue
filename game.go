package gogue

import (
  "errors"
)

type Game struct {
  Map
  Player
}

func (g Game) FollowPath(steps... rune) (Game, error) {
  if len(steps) == 0 {
    return g, nil
  }

  var err error

  switch steps[0] {
  case 'n': g, err = g.MoveNorth()
  case 's': g, err = g.MoveSouth()
  case 'e': g, err = g.MoveEast()
  case 'w': g, err = g.MoveWest()
  }

  if (err != nil) {
    return g, err
  }

  return g.FollowPath(steps[1:]...)
}

func (g Game) IsOver() bool {
  return g.Player.Position == g.Goal.Position
}

func (g Game) MoveNorth() (Game, error) {
  if g.Map.Get(g.Player.X, g.Player.Y - 1).IsWalkable() {
    g.Player.Y -= 1
    return g, nil
  } else {
    err := errors.New("Cannot move")
    return g, err
  }
}

func (g Game) MoveSouth() (Game, error) {
  if g.Map.Get(g.Player.X, g.Player.Y + 1).IsWalkable() {
    g.Player.Y += 1
    return g, nil
  } else {
    err := errors.New("Cannot move")
    return g, err
  }
}

func (g Game) MoveEast() (Game, error) {
  if g.Map.Get(g.Player.X + 1, g.Player.Y).IsWalkable() {
    g.Player.X += 1
    return g, nil
  } else {
    err := errors.New("Cannot move")
    return g, err
  }
}

func (g Game) MoveWest() (Game, error) {
  if g.Map.Get(g.Player.X - 1, g.Player.Y).IsWalkable() {
    g.Player.X -= 1
    return g, nil
  } else {
    err := errors.New("Cannot move")
    return g, err
  }
}
