package server_test

import (
	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	. "github.com/luan/gogue/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		mmap      *gogue.Map
		broadcast chan protocol.Packet
		client    *Client
	)

	BeforeEach(func() {
		broadcast = make(chan protocol.Packet)
		mmap = gogue.NewMap("...")
		client = NewClient(mmap, broadcast)
		go client.Run()
	})

	Describe("when connecting", func() {
		It("broadcasts its location", func(done Done) {
			Expect(<-broadcast).To(Equal(protocol.Creature{protocol.Position{1, 1, 0}}))
			close(done)
		})

		It("sends out the visible map", func(done Done) {
			Expect(<-client.Outgoing).To(Equal(protocol.MapPortion{"...\n"}))
			close(done)
		})
	})
})
