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
		mmap = gogue.NewMap("...\n...\n...")
		client = NewClient(mmap, broadcast)
		go client.Run()
	})

	Describe("CreaturePacket", func() {
		It("has the players UUID", func() {
			Expect(client.CreaturePacket().UUID).To(Equal(client.Player.UUID))
		})

		It("has the players Position", func() {
			Expect(client.CreaturePacket().Position).To(
				Equal(protocol.Position(client.Player.Position)))
		})
	})

	Describe("when connecting", func() {
		It("broadcasts its location", func(done Done) {
			Eventually(func() protocol.Packet {
				return <-broadcast
			}).Should(Equal(client.CreaturePacket()))
			close(done)
		})

		It("sends out the visible map", func(done Done) {
			Eventually(func() protocol.Packet {
				return <-client.Outgoing
			}).Should(Equal(protocol.MapPortion{"...\n...\n...\n"}))
			close(done)
			close(done)
		})
	})
})
