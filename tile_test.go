package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tile", func() {
	Describe("IsWalkable", func() {
		It("explicitly walkable tiles are walkable", func() {
			tile := Tile{Properties: map[string]string{"walkable": "true"}}
			Expect(tile.IsWalkable()).To(BeTrue())
		})

		It("explicitly not walkable tiles aren't walkable", func() {
			tile := Tile{Properties: map[string]string{"walkable": "false"}}
			Expect(tile.IsWalkable()).To(BeFalse())
		})

		It("any other tile is not walkable", func() {
			Expect(Tile{}.IsWalkable()).To(BeTrue())
		})
	})

	Describe("PositionModifier", func() {
		It("gives back the X modifier", func() {
			tile := Tile{Properties: map[string]string{"changeX": "5"}}
			Expect(tile.PositionModifier()).To(Equal(Position{5, 0, 0}))
		})

		It("gives back the Y modifier", func() {
			tile := Tile{Properties: map[string]string{"changeY": "-2"}}
			Expect(tile.PositionModifier()).To(Equal(Position{0, -2, 0}))
		})

		It("gives back the Z modifier", func() {
			tile := Tile{Properties: map[string]string{"changeZ": "2"}}
			Expect(tile.PositionModifier()).To(Equal(Position{0, 0, 2}))
		})

		It("gives back the all modifiers", func() {
			tile := Tile{Properties: map[string]string{"changeX": "-1", "changeY": "-1", "changeZ": "-1"}}
			Expect(tile.PositionModifier()).To(Equal(Position{-1, -1, -1}))
		})
	})
})
