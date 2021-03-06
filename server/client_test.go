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
		broadcast chan protocol.Packet
		client    *Client
		mmap      *gogue.Map
		quit      chan bool
	)

	BeforeEach(func() {
		broadcast = make(chan protocol.Packet)
		quit = make(chan bool)
		mmap = gogue.NewMap("../assets/map-tiled.json")
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
			<-client.Outgoing
			<-client.Outgoing
			Eventually(func() protocol.Packet {
				return <-broadcast
			}).Should(Equal(client.CreaturePacket()))
			close(done)
		})

		It("sends out the visible map", func(done Done) {
			Eventually(func() protocol.Packet {
				p := <-client.Outgoing
				return p
			}).Should(Equal(protocol.Map{*mmap}))
			close(done)
		})

		It("sends out the player uuid", func(done Done) {
			Eventually(func() protocol.Packet {
				p := <-client.Outgoing
				return p
			}).Should(Equal(protocol.Player{client.UUID}))
			close(done)
		})
	})

	Context("after connected", func() {
		BeforeEach(func(done Done) {
			<-client.Outgoing
			<-client.Outgoing
			<-broadcast
			close(done)
		})

		Describe("logging out", func() {
			BeforeEach(func(done Done) {
				client.Incoming <- protocol.Quit{}
				close(done)
			})

			It("closes the quit channel", func(done Done) {
				Eventually(client.Quit).Should(BeClosed())
				close(done)
			})
		})

		Describe("walking", func() {
			It("broadcasts its new location", func(done Done) {
				client.Incoming <- protocol.WalkNorth
				Expect(<-broadcast).To(Equal(protocol.Creature{
					UUID:     client.UUID,
					Position: protocol.Position{X: 1, Y: 0, Z: 0},
				}))
				close(done)
			})
		})
	})
})
