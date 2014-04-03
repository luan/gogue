package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	Describe("NewMap", func() {
		It("creates a map with the content passed in", func() {
			newMap := NewMap(`
      ..#
      #.*`)
			Expect(newMap.Height).To(Equal(2))
			Expect(newMap.Width).To(Equal(3))
			Expect(newMap.Depth).To(Equal(1))
			Expect(newMap.Get(Position{X: 0, Y: 0})).To(Equal(Tile('.')))
			Expect(newMap.Get(Position{X: 2, Y: 1})).To(Equal(Tile('*')))
		})

		It("has multiple floors", func() {
			newMap := NewMap(`
      ..#
      #.>
      `, `
      #.#
      >.<`, `
      ..*
      <##`)
			Expect(newMap.Height).To(Equal(2))
			Expect(newMap.Width).To(Equal(3))
			Expect(newMap.Depth).To(Equal(3))
			Expect(newMap.Get(Position{X: 0, Y: 0, Z: 0})).To(Equal(Tile('.')))
			Expect(newMap.Get(Position{X: 2, Y: 1, Z: 0})).To(Equal(Tile('>')))
			Expect(newMap.Get(Position{X: 2, Y: 0, Z: 2})).To(Equal(Tile('*')))
		})
	})
})
