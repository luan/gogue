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
				Eventually(client.Receive).Should(BeAssignableToTypeOf(protocol.MapPortion{}))
				close(done)
			})

			It("broadcasts its presence to all connected clients", func(done Done) {
				client1 := fakes.NewClient()
				client1.Connect(listener)
				client2 := fakes.NewClient()
				client2.Connect(listener)

				Eventually(client1.Receive).Should(
					Equal(protocol.Creature{protocol.Position{1, 1, 0}}))
				close(done)
			})
		})
	})
})
