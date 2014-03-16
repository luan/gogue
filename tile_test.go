package gogue_test

import (
	. "github.com/luan/gogue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tile", func() {
  Describe("IsWalkable", func() {
    It("Knows that '.' are walkable", func() {
      Expect(Tile('.').IsWalkable()).To(BeTrue())
    })

    It("Knows that '*' are walkable", func() {
      Expect(Tile('*').IsWalkable()).To(BeTrue())
    })

    It("And everything else is not", func() {
      Expect(Tile('#').IsWalkable()).To(BeFalse())
      Expect(Tile('$').IsWalkable()).To(BeFalse())
      Expect(Tile('-').IsWalkable()).To(BeFalse())
      Expect(Tile('@').IsWalkable()).To(BeFalse())
      Expect(Tile('3').IsWalkable()).To(BeFalse())
    })
  })
})