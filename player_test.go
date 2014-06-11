package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Player", func() {
	var (
		player *Player
		mmap   *Map
		err    error
	)

	JustBeforeEach(func() {
		mmap = NewMap("assets/map-tiled.json")
		player = NewPlayer("me", mmap, Position{4, 5, 0})
	})

	Describe("Walking", func() {
		It("cannot move north because there is a wall", func() {
			err = player.MoveNorth()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("cannot move"))
			Expect(player.X).To(Equal(4))
			Expect(player.Y).To(Equal(5))
		})

		It("Can move south", func() {
			err = player.MoveSouth()
			Expect(err).NotTo(HaveOccurred())
			Expect(player.X).To(Equal(4))
			Expect(player.Y).To(Equal(6))
		})

		It("Can move east", func() {
			err = player.MoveEast()
			Expect(err).NotTo(HaveOccurred())
			Expect(player.X).To(Equal(5))
			Expect(player.Y).To(Equal(5))
		})

		It("cannot move west because there is a wall", func() {
			err = player.MoveWest()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("cannot move"))
			Expect(player.X).To(Equal(4))
			Expect(player.Y).To(Equal(5))
		})

		Context("when next tile is position modifier", func() {
			It("goes down when move towards the stairs", func() {
				player.X = 9
				player.Y = 7
				err = player.MoveEast()
				Expect(err).NotTo(HaveOccurred())
				Expect(player.X).To(Equal(9))
				Expect(player.Y).To(Equal(5))
				Expect(player.Z).To(Equal(1))
			})
		})
	})
})
