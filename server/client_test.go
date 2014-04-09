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
		})
	})

	Describe("walking", func() {
		It("updates its location", func(done Done) {
			client.Incoming <- protocol.Walk{protocol.North}
			Eventually(func() protocol.Position {
				return client.CreaturePacket().Position
			}).Should(Equal(protocol.Position{1, 0, 0}))
			close(done)
		})

		It("broadcasts its new location", func(done Done) {
			client.Incoming <- protocol.Walk{protocol.East}
			Eventually(func() protocol.Position {
				return client.CreaturePacket().Position
			}).Should(Equal(protocol.Position{2, 1, 0}))

			Eventually(func() protocol.Packet {
				return <-broadcast
			}).Should(Equal(client.CreaturePacket()))
			close(done)
		})

		It("receives its new location", func(done Done) {
			client.Incoming <- protocol.Walk{protocol.West}
			Eventually(func() protocol.Position {
				return client.CreaturePacket().Position
			}).Should(Equal(protocol.Position{0, 1, 0}))

			Eventually(func() protocol.Packet {
				return <-client.Outgoing
			}).Should(Equal(client.CreaturePacket()))
			close(done)
		})
	})
})
