package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Game", func() {
	var (
		player   Player
		game     Game
		gameMap  Map
		mapInput []string
		err      error
	)
	BeforeEach(func() {
		mapInput = []string{`
		######
		#....#
		#..#*#
		#....#
		######
		`}
	})

	JustBeforeEach(func() {
		player = Player{Position{X: 1, Y: 1}}
		gameMap, _ = NewMap(mapInput...)
		game = Game{gameMap, player}
	})

	It("ends when the player reaches the goal", func() {
		game, _ = game.MoveEast()
		game, _ = game.MoveEast()
		game, _ = game.MoveEast()
		Expect(game.Player.Position.X).To(Equal(4))
		Expect(game.Player.Position.Y).To(Equal(1))
		Expect(game.IsOver()).To(BeFalse())

		game, _ = game.MoveSouth()
		Expect(game.Player.Position.X).To(Equal(4))
		Expect(game.Player.Position.Y).To(Equal(2))
		Expect(game.IsOver()).To(BeTrue())
	})

	Describe("Walking", func() {
		Context("When all squares around the player are walkable", func() {
			BeforeEach(func() {
				mapInput = []string{`
				...
				...
				...
				`}
			})

			It("Can move north", func() {
				game, err = game.MoveNorth()
				Expect(err).To(BeNil())
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(0))
			})

			It("Can move south", func() {
				game, err = game.MoveSouth()
				Expect(err).To(BeNil())
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(2))
			})

			It("Can move east", func() {
				game, err = game.MoveEast()
				Expect(err).To(BeNil())
				Expect(game.Player.X).To(Equal(2))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("Can move west", func() {
				game, err = game.MoveWest()
				Expect(err).To(BeNil())
				Expect(game.Player.X).To(Equal(0))
				Expect(game.Player.Y).To(Equal(1))
			})
		})

		Context("When no squares around the player are walkable", func() {
			BeforeEach(func() {
				mapInput = []string{`
				###
				#.#
				###
				`}
			})

			It("cannot move north", func() {
				game, err = game.MoveNorth()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("cannot move south", func() {
				game, err = game.MoveSouth()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("cannot move east", func() {
				game, err = game.MoveEast()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("cannot move west", func() {
				game, err = game.MoveWest()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})
		})

		Context("when next tile is floor changer", func() {
			BeforeEach(func() {
				mapInput = []string{`
				...
				>..
				`, `
				...
				<.*
				`}
			})

			It("goes down when move towards the stairs", func() {
				game, err = game.MoveWest()
				Expect(err).To(BeNil())
				Expect(game.Player.X).To(Equal(0))
				Expect(game.Player.Y).To(Equal(1))
				Expect(game.Player.Z).To(Equal(1))
			})

			It("can go back up", func() {
				game, _ = game.MoveWest()
				game, _ = game.MoveEast()
				game, err = game.MoveWest()
				Expect(err).To(BeNil())
				Expect(game.Player.X).To(Equal(0))
				Expect(game.Player.Y).To(Equal(1))
				Expect(game.Player.Z).To(Equal(0))
			})
		})
	})
})
