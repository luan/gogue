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
      #.*
      `)
			Expect(err).To(BeNil())
			Expect(newMap.Height).To(Equal(2))
			Expect(newMap.Width).To(Equal(3))
			Expect(newMap.Get(0, 0)).To(Equal(Tile('.')))
			Expect(newMap.Get(2, 1)).To(Equal(Tile('*')))
		})

		It("errors if there is no goal", func() {
			_, err := NewMap(`.`)
			Expect(err.Error()).To(Equal("Map requires a Goal(*)"))
		})

		It("knows where the goal is", func() {
			newMap, _ := NewMap(`....*`)
			Expect(newMap.Goal.X).To(Equal(4))
			Expect(newMap.Goal.Y).To(Equal(0))
		})
	})
})
