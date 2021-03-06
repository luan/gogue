package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Position", func() {
	Describe("North", func() {
		It("returns a position one to the north of the current", func() {
			pos := Position{X: 2, Y: 1}
			Expect(pos.North()).To(Equal(Position{X: 2, Y: 0}))
		})
	})

	Describe("South", func() {
		It("returns a position one to the south of the current", func() {
			pos := Position{X: 1, Y: 2}
			Expect(pos.South()).To(Equal(Position{X: 1, Y: 3}))
		})
	})

	Describe("East", func() {
		It("returns a position one to the east of the current", func() {
			pos := Position{X: 1, Y: 2}
			Expect(pos.East()).To(Equal(Position{X: 2, Y: 2}))
		})
	})

	Describe("West", func() {
		It("returns a position one to the west of the current", func() {
			pos := Position{X: 2, Y: 2}
			Expect(pos.West()).To(Equal(Position{X: 1, Y: 2}))
		})
	})

	Describe("Up", func() {
		It("changes the position to the upper floor of the current", func() {
			pos := Position{X: 2, Y: 2, Z: 2}
			Expect(pos.Up()).To(Equal(Position{X: 2, Y: 2, Z: 3}))
		})
	})

	Describe("Down", func() {
		It("changes the position to the lower floor of the current", func() {
			pos := Position{X: 2, Y: 2, Z: 2}
			Expect(pos.Down()).To(Equal(Position{X: 2, Y: 2, Z: 1}))
		})
	})

	Describe("Add", func() {
		It("adds itself to another position", func() {
			pos := Position{X: 2, Y: 2, Z: 2}
			Expect(pos.Add(Position{1, 1, 1})).To(Equal(Position{X: 3, Y: 3, Z: 3}))
		})
	})
})
