package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	Describe("NewMap", func() {
		It("creates a map with the content passed in", func() {
			newMap, err := NewMap(`
      ..#
      #.*`)
			Expect(err).To(BeNil())
			Expect(newMap.Height).To(Equal(2))
			Expect(newMap.Width).To(Equal(3))
			Expect(newMap.Depth).To(Equal(1))
			Expect(newMap.Get(Position{X: 0, Y: 0})).To(Equal(Tile('.')))
			Expect(newMap.Get(Position{X: 2, Y: 1})).To(Equal(Tile('*')))
		})

		It("has multiple floors", func() {
			newMap, err := NewMap(`
      ..#
      #.>
      `, `
      #.#
      >.<`, `
      ..*
      <##`)
			Expect(err).To(BeNil())
			Expect(newMap.Height).To(Equal(2))
			Expect(newMap.Width).To(Equal(3))
			Expect(newMap.Depth).To(Equal(3))
			Expect(newMap.Get(Position{X: 0, Y: 0, Z: 0})).To(Equal(Tile('.')))
			Expect(newMap.Get(Position{X: 2, Y: 1, Z: 0})).To(Equal(Tile('>')))
			Expect(newMap.Get(Position{X: 2, Y: 0, Z: 2})).To(Equal(Tile('*')))
		})

		It("errors if there is no goal", func() {
			_, err := NewMap(`.`)
			Expect(err.Error()).To(Equal("map requires a Goal(*)"))
		})

		It("knows where the goal is", func() {
			newMap, _ := NewMap(`....>`, `..*..<`)
			Expect(newMap.Goal.X).To(Equal(2))
			Expect(newMap.Goal.Y).To(Equal(0))
			Expect(newMap.Goal.Z).To(Equal(1))
		})
	})
})
