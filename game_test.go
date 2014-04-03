package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Game", func() {
	var (
		game     *Game
		mmap     *Map
		mapInput []string
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
		mmap = NewMap(mapInput...)
		game = NewGame(mmap)
	})

	Describe("AddPlayer", func() {
		It("adds a player to the game", func() {
			game.AddPlayer("1", Position{1, 1, 0})
			Expect(game.Players["1"].X).To(Equal(1))
		})

		It("returns the created player", func() {
			player := game.AddPlayer("4", Position{1, 1, 0})
			Expect(player.Guid).To(Equal("4"))
		})
	})
})
