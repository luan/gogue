package server_test

import (
	"time"

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

		Context("when a client disconnects", func() {
			It("forgets about the client and stops broadcasting to it", func(done Done) {
				rcl1 := fakes.NewClient()
				rcl1.Connect(listener)
				rcl2 := fakes.NewClient()
				rcl2.Connect(listener)

				rcl1.Outgoing <- protocol.Quit{}

				Eventually(func() protocol.Packet {
					return <-rcl2.Incoming
				}).Should(BeAssignableToTypeOf(protocol.RemoveCreature{}))

				rcl2.Outgoing <- protocol.WalkNorth

			noincoming:
				for {
					select {
					case p := <-rcl1.Incoming:
						if cp, ok := p.(protocol.Creature); ok {
							Expect(cp.Y).NotTo(Equal(0))
						}
					default:
						break noincoming
					}
					time.Sleep(100)
				}

				close(done)
			})
		})
	})
})
