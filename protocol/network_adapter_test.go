package protocol_test

import (
	. "github.com/luan/gogue/protocol"
	"github.com/luan/gogue/test/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworkAdapter", func() {
	var (
		na       *NetworkAdapter
		conn     *fakes.Conn
		incoming chan Packet
		outgoing chan Packet
		quit     chan bool
	)

	BeforeEach(func() {
		incoming = make(chan Packet)
		outgoing = make(chan Packet)
		quit = make(chan bool)
		conn = fakes.NewConn()
		na = NewNetworkAdapter(incoming, outgoing, quit, conn)
		na.Listen()
	})

	Describe("listening", func() {
		It("passes packets from the wire to the incoming channel", func(done Done) {
			p := Quit{}
			conn.Send(p)
			Expect(<-incoming).To(Equal(p))
			close(done)
		})

		It("passes along packets from the outgoing channel to the wire", func(done Done) {
			p := RemoveCreature{"abc"}
			outgoing <- p
			Expect(conn.Receive()).To(Equal(p))
			close(done)
		})

		It("stops everything upon quitting", func(done Done) {
			quit <- true
			Eventually(conn.Closed).Should(BeTrue())
			close(done)
		})
	})
})
