package gogue

import "errors"

type Player struct { Position }

type Game struct {
  Map
  Player
}

func (game Game) FollowPath(steps... rune) (Game, error) {
  if len(steps) == 0 {
    return game, nil
  }

  var err error

  switch steps[0] {
  case 'n': game, err = game.MoveNorth()
  case 's': game, err = game.MoveSouth()
  case 'e': game, err = game.MoveEast()
  case 'w': game, err = game.MoveWest()
  }

  if (err != nil) {
    return game, err
  }

  return game.FollowPath(steps[1:]...)
}

func (game Game) IsOver() bool {
  return game.Player.Position == game.Goal.Position
}

func (game Game) MoveNorth() (Game, error) {
  if game.Map.Get(game.Player.X, game.Player.Y - 1).IsWalkable() {
    game.Player.Y -= 1
    return game, nil
  } else {
    err := errors.New("Cannot move")
    return game, err
  }
}

func (game Game) MoveSouth() (Game, error) {
  if game.Map.Get(game.Player.X, game.Player.Y + 1).IsWalkable() {
    game.Player.Y += 1
    return game, nil
  } else {
    err := errors.New("Cannot move")
    return game, err
  }
}

func (game Game) MoveEast() (Game, error) {
  if game.Map.Get(game.Player.X + 1, game.Player.Y).IsWalkable() {
    game.Player.X += 1
    return game, nil
  } else {
    err := errors.New("Cannot move")
    return game, err
  }
}

func (game Game) MoveWest() (Game, error) {
  if game.Map.Get(game.Player.X - 1, game.Player.Y).IsWalkable() {
    game.Player.X -= 1
    return game, nil
  } else {
    err := errors.New("Cannot move")
    return game, err
  }
}
