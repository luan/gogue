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
		mmap = gogue.NewMap("../assets/map-tiled.json")
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
				rcl.Outgoing <- protocol.WalkEast
				for i := 0; i < 3; i++ {
					rcl := fakes.NewClient()
					rcl.Connect(listener)
					time.Sleep(100)
					remoteClients = append(remoteClients, rcl)
				}
				close(done)
			})

			It("sends the new client map for the current floor", func() {
				rcl := fakes.NewClient()
				rcl.Connect(listener)
				Eventually(rcl.ReceivedPackets).Should(ContainElement(protocol.Map{*mmap}))
			})

			It("sends the new client info about other players at the current floor", func() {
				rcl := fakes.NewClient()
				rcl.Connect(listener)

				for _, cl := range gs.Clients() {
					Eventually(rcl.ReceivedPackets).Should(ContainElement(cl.CreaturePacket()))
				}
			})
		})
	})
})
