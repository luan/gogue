package server_test

import (
	"strings"

	"github.com/luan/gogue"
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
		game     *gogue.Game
	)

	BeforeEach(func() {
		listener = fakes.NewListener()
		game = gogue.NewGame(gogue.NewMap(`
		###############
		#.............#
		#.............#
		#.............#
		#.............#
		#.............#
		###############
		`))
		gs = NewGameServer(game, listener)
		go gs.WaitForClients()
	})

	Describe("client connections", func() {
		Context("when a new client connects", func() {
			It("adds the client to the client list", func() {
				client := fakes.NewClient()
				client.Connect(listener)

				Eventually(func() []*Client {
					return gs.Clients
				}).Should(HaveLen(1))
			})

			It("sends the visible map to the connected client", func() {
				client := fakes.NewClient()
				client.Connect(listener)

				var packet protocol.MapPortion

				Eventually(func() (ok bool) {
					if p, err := client.Receive(); err == nil {
						packet, ok = p.(protocol.MapPortion)
					}
					return
				}).Should(BeTrue())

				Expect(packet).To(Equal(protocol.MapPortion{strings.Replace(`###############
				#.............#
				#.............#
				#.............#
				#.............#
				#.............#
				###############
				`, "\t", "", -1)}))
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
