package server_test

import (
	"github.com/luan/gogue"
	"github.com/luan/gogue/protocol"
	. "github.com/luan/gogue/server"
	"github.com/luan/gogue/test/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GameServer", func() {
	var (
		listener *fakes.Listener
		gs       *GameServer
		mmap     *gogue.Map
	)

	BeforeEach(func() {
		listener = fakes.NewListener()
		mmap = gogue.NewMap("#...#\n#...#")
		gs = NewGameServer(mmap, listener)
		go gs.Run()
	})

	Describe("client connections", func() {
		Context("when a new client connects", func() {
			It("sends the visible map to the connected client", func(done Done) {
				client := fakes.NewClient()
				client.Connect(listener)
				Eventually(client.Receive).Should(
					BeAssignableToTypeOf(protocol.MapPortion{}))
				close(done)
			})

			It("broadcasts its presence to all connected clients", func(done Done) {
				client1 := fakes.NewClient()
				client1.Connect(listener)
				client2 := fakes.NewClient()
				client2.Connect(listener)

				Eventually(client1.Receive).Should(
					BeAssignableToTypeOf(protocol.Creature{}))
				close(done)
			})
		})

		Context("when a client disconnects", func() {
			It("forgets about the client and stops broadcasting to it", func(done Done) {
				client1 := fakes.NewClient()
				client1.Connect(listener)
				client2 := fakes.NewClient()
				client2.Connect(listener)

				Eventually(client1.Receive).Should(BeAssignableToTypeOf(protocol.Creature{}))
				Eventually(client2.Receive).Should(BeAssignableToTypeOf(protocol.Creature{}))

				client1.Send(protocol.Quit{})
				Eventually(client2.Receive).Should(BeAssignableToTypeOf(protocol.RemoveCreature{}))

				client2.Send(protocol.Walk{protocol.East})
				Eventually(client1.Receive).Should(BeNil())
				close(done)
			})
		})
	})
})
