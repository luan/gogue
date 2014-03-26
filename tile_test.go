package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tile", func() {
	Describe("IsWalkable", func() {
		It("knows that '.' are walkable", func() {
			Expect(Tile('.').IsWalkable()).To(BeTrue())
		})

		It("knows that '*' are walkable", func() {
			Expect(Tile('*').IsWalkable()).To(BeTrue())
		})

		It("knows that floor changers are walkable", func() {
			Expect(Tile('>').IsWalkable()).To(BeTrue())
			Expect(Tile('<').IsWalkable()).To(BeTrue())
		})

		It("and everything else is not", func() {
			Expect(Tile('#').IsWalkable()).To(BeFalse())
			Expect(Tile('$').IsWalkable()).To(BeFalse())
			Expect(Tile('-').IsWalkable()).To(BeFalse())
			Expect(Tile('@').IsWalkable()).To(BeFalse())
			Expect(Tile('3').IsWalkable()).To(BeFalse())
		})
	})

	Describe("ChangeFloor", func() {
		It("goes down on '>' tiles", func() {
			Expect(Tile('>').ChangeFloor()).To(Equal("down"))
		})

		It("goes up on '<' tiles", func() {
			Expect(Tile('<').ChangeFloor()).To(Equal("up"))
		})

		It("stays still on other tiles", func() {
			Expect(Tile('.').ChangeFloor()).To(Equal(""))
		})
	})
})
