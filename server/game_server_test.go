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
		mmap = gogue.NewMap("#...#\n#.>.#", "#...#\n.<.#")
		gs = NewGameServer(mmap, listener)
		go gs.Run()
	})

	Describe("client connections", func() {
		Context("there are other players connected", func() {
			var (
				remoteClients []*fakes.Client
			)

			BeforeEach(func(done Done) {
				rcl := fakes.NewClient()
				rcl.Connect(listener)
				remoteClients = append(remoteClients, rcl)
				rcl.Outgoing <- protocol.WalkEast
				for i := 0; i < 3; i++ {
					rcl := fakes.NewClient()
					rcl.Connect(listener)
					remoteClients = append(remoteClients, rcl)
				}
				for _, rcl := range remoteClients {
				inner:
					for {
						select {
						case <-rcl.Incoming:
						default:
							break inner
						}
					}
				}
				close(done)
			})

			It("sends the new client the map for the current floor", func(done Done) {
				rcl := fakes.NewClient()
				rcl.Connect(listener)
				Eventually(func() protocol.Packet {
					return <-rcl.Incoming
				}).Should(Equal(protocol.MapPortion{Data: "#...#\n#.>.#\n"}))
				close(done)
			})

			It("sends the new client info about other players on the same floor", func(done Done) {
				rcl := fakes.NewClient()
				rcl.Connect(listener)
				packets := []protocol.Packet{}
				<-rcl.Incoming // map packet
				for i := 0; i < 4; i++ {
					packets = append(packets, <-rcl.Incoming)
				}

				for _, cl := range gs.Clients() {
					if cl.Player.Z == 0 {
						Expect(packets).To(ContainElement(cl.CreaturePacket()))
					} else {
						Expect(packets).NotTo(ContainElement(cl.CreaturePacket()))
					}
				}
				close(done)
			})

			It("sends the new clients info to all other players on the same floor", func(done Done) {
				rcl := fakes.NewClient()
				rcl.Connect(listener)
				for _, rocl := range remoteClients {
					Eventually(func() protocol.Packet {
						return <-rocl.Incoming
					}).Should(BeAssignableToTypeOf(protocol.Creature{}))
				}
				close(done)
			})
		})
	})
})
