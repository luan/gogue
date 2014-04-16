package protocol_test

import (
	"encoding/gob"
	"net"
	. "github.com/luan/gogue/protocol"
	"github.com/luan/gogue/test/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworkAdapter", func() {
	var (
		na           *NetworkAdapter
		serverBuffer net.Conn
		clientBuffer net.Conn
		conn         *fakes.Conn
		incoming     chan Packet
		outgoing     chan Packet
		quit         chan bool
	)

	BeforeEach(func() {
		incoming = make(chan Packet)
		outgoing = make(chan Packet)
		quit = make(chan bool)
		serverBuffer, clientBuffer = net.Pipe()
		conn = fakes.NewConn(serverBuffer)
		na = NewNetworkAdapter(incoming, outgoing, quit, conn)
		na.Listen()
	})

	Describe("listening", func() {
		It("passes packets from the wire to the incoming channel", func(done Done) {
			p := Quit{}
			enc := gob.NewEncoder(clientBuffer)
			enc.Encode(&p)
			Expect(<-incoming).To(Equal(p))
			close(done)
		})

		It("passes along packets from the outgoing channel to the wire", func(done Done) {
			p := RemoveCreature{"abc"}
			outgoing <- p

			var q Packet
			dec := gob.NewDecoder(clientBuffer)
			dec.Decode(&q)
			Expect(q).To(Equal(p))
			close(done)
		})

		It("stops when the connection is closed", func(done Done) {
			conn.Close()
			Eventually(func() Packet {
				return <-incoming
			}).Should(Equal(Quit{}))
			close(done)
		})

		It("stops when the connection is closed", func(done Done) {
			conn.Close()
			outgoing <- Creature{Position: Position{0, 0, 0}}
			Eventually(func() Packet {
				return <-incoming
			}).Should(Equal(Quit{}))
			close(done)
		})

		It("stops everything upon quitting", func(done Done) {
			quit <- true
			Eventually(conn.Closed).Should(BeTrue())
			close(done)
		})
	})
})
