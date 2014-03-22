package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Game", func() {
	var (
		player    Player
		game      Game
		gameMap   Map
		mapString string
		err       error
	)
	BeforeEach(func() {
		mapString = `
		##########
		#........#
		####.#####
		#....###*#
		#.######.#
		#........#
		##########
		`
	})

	JustBeforeEach(func() {
		player = Player{Position{1, 1}}
		gameMap, _ = NewMap(mapString)
		game = Game{gameMap, player}
	})

	It("ends when the player reaches the goal", func() {
		game, err = game.FollowPath('e', 'e', 'e', 's', 's')
		Expect(game.Player.Position.X).To(Equal(4))
		Expect(game.Player.Position.Y).To(Equal(3))
		Expect(game.IsOver()).To(BeFalse())

		game, err = game.FollowPath('w', 'w', 'w', 's', 's')
		Expect(game.Player.Position.X).To(Equal(1))
		Expect(game.Player.Position.Y).To(Equal(5))
		Expect(game.IsOver()).To(BeFalse())

		game, err = game.FollowPath('e', 'e', 'e', 'e', 'e', 'e', 'e', 'n', 'n')
		Expect(game.IsOver()).To(BeTrue())
		Expect(game.Player.Position.X).To(Equal(8))
		Expect(game.Player.Position.Y).To(Equal(3))
	})

	Describe("Walking", func() {
		Context("When all squares around the player are walkable", func() {
			BeforeEach(func() {
				mapString = `
				...
				...
				...
				`
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
				mapString = `
				###
				#.#
				###
				`
			})

			It("Cannot move north", func() {
				game, err = game.MoveNorth()
				Expect(err.Error()).To(Equal("Cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("Cannot move south", func() {
				game, err = game.MoveSouth()
				Expect(err.Error()).To(Equal("Cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("Cannot move east", func() {
				game, err = game.MoveEast()
				Expect(err.Error()).To(Equal("Cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})

			It("Cannot move west", func() {
				game, err = game.MoveWest()
				Expect(err.Error()).To(Equal("Cannot move"))
				Expect(game.Player.X).To(Equal(1))
				Expect(game.Player.Y).To(Equal(1))
			})
		})
	})
})
