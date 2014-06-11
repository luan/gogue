package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	Describe("NewMap", func() {
		It("creates a map from the file passed in", func() {
			newMap := NewMap("assets/map-tiled.json")
			var tile Tile
			tile, _ = newMap.Get(Position{X: 10, Y: 12, Z: 0})
			Expect(tile.Tiles).To(Equal([]int{7, 0}))
			Expect(tile.PositionModifier()).To(Equal(Position{1, 2, -1}))
			Expect(tile.IsWalkable()).To(BeTrue())

			tile, _ = newMap.Get(Position{X: 10, Y: 13, Z: 0})
			Expect(tile.Tiles).To(Equal([]int{5, 2}))
			Expect(tile.PositionModifier()).To(Equal(Position{0, 0, 0}))
			Expect(tile.IsWalkable()).To(BeFalse())

			tile, _ = newMap.Get(Position{X: 11, Y: 13, Z: -1})
			Expect(tile.Tiles).To(Equal([]int{6, 0}))
			Expect(tile.PositionModifier()).To(Equal(Position{-1, -2, 1}))
			Expect(tile.IsWalkable()).To(BeTrue())

			tile, _ = newMap.Get(Position{X: 9, Y: 6, Z: 1})
			Expect(tile.Tiles).To(Equal([]int{7, 0}))
			Expect(tile.PositionModifier()).To(Equal(Position{1, 2, -1}))
			Expect(tile.IsWalkable()).To(BeTrue())

			tile, _ = newMap.Get(Position{X: 11, Y: 7, Z: 1})
			Expect(tile.Tiles).To(Equal([]int{5, 3}))
			Expect(tile.IsWalkable()).To(BeFalse())
		})
	})
})
