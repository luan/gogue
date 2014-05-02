package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Player", func() {
	var (
		player   *Player
		mmap     *Map
		mapInput []string
		err      error
	)

	JustBeforeEach(func() {
		mmap = NewMap(mapInput...)
		player = NewPlayer("me", mmap, Position{1, 1, 0})
	})

	Describe("MapSight", func() {
		BeforeEach(func() {
			mapInput = []string{`
				...
				.#.
				...
				`}
		})

		It("returns the visible portion of the map", func() {
			sight := player.MapSight()
			Expect(sight).To(MatchRegexp(strings.Replace(`...
			.#.
			...
			`, "\t", "", -1)))
		})
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
				err = player.MoveNorth()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(1))
				Expect(player.Y).To(Equal(0))
			})

			It("Can move south", func() {
				err = player.MoveSouth()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(1))
				Expect(player.Y).To(Equal(2))
			})

			It("Can move east", func() {
				err = player.MoveEast()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(2))
				Expect(player.Y).To(Equal(1))
			})

			It("Can move west", func() {
				err = player.MoveWest()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(0))
				Expect(player.Y).To(Equal(1))
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
				err = player.MoveNorth()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(player.X).To(Equal(1))
				Expect(player.Y).To(Equal(1))
			})

			It("cannot move south", func() {
				err = player.MoveSouth()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(player.X).To(Equal(1))
				Expect(player.Y).To(Equal(1))
			})

			It("cannot move east", func() {
				err = player.MoveEast()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(player.X).To(Equal(1))
				Expect(player.Y).To(Equal(1))
			})

			It("cannot move west", func() {
				err = player.MoveWest()
				Expect(err.Error()).To(Equal("cannot move"))
				Expect(player.X).To(Equal(1))
				Expect(player.Y).To(Equal(1))
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
				err = player.MoveWest()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(0))
				Expect(player.Y).To(Equal(1))
				Expect(player.Z).To(Equal(1))
			})

			It("can go back up", func() {
				player.MoveWest()
				player.MoveEast()
				err = player.MoveWest()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(0))
				Expect(player.Y).To(Equal(1))
				Expect(player.Z).To(Equal(0))
			})
		})
	})
})
