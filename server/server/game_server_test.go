package server_test

import (
	"github.com/luan/gogue/protocol"
	. "github.com/luan/gogue/server/server"
	"github.com/luan/gogue/test/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GameServer", func() {
	var (
		listener *fakes.Listener
		gs       *GameServer
	)

	BeforeEach(func() {
		listener = fakes.NewListener()
		gs = NewGameServer(listener)
		go gs.WaitForClients()
	})

	Describe("client connections", func() {
		Context("when a new client connects", func() {
			It("adds the client to the client list", func() {
				client := fakes.NewClient()
				client.Connect(listener)

				Eventually(func() []Client {
					return gs.Clients
				}).Should(ContainElement(client))
			})

			It("broadcasts its presence to all connected clients", func() {
				client1 := fakes.NewClient()
				client1.Connect(listener)
				client2 := fakes.NewClient()
				client2.Connect(listener)

				var packet protocol.Creature

				Eventually(func() (ok bool) {
					if p, err := client1.Receive(); err == nil {
						packet, ok = p.(protocol.Creature)
					}
					return
				}).Should(BeTrue())

				Expect(packet).To(Equal(protocol.Creature{protocol.Position{1, 1, 0}}))
			})
		})
	})
})
