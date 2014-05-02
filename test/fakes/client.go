package fakes

import (
	"net"
	"sync"

	"github.com/luan/gogue/protocol"
)

type Client struct {
	remote          *Conn
	local           *Conn
	Outgoing        chan protocol.Packet
	Incoming        chan protocol.Packet
	Quit            chan bool
	receivedPackets []protocol.Packet
	packetsMutex    sync.Mutex
}

func NewClient() *Client {
	serverBuffer, clientBuffer := net.Pipe()
	return &Client{
		remote:   NewConn(serverBuffer),
		local:    NewConn(clientBuffer),
		Outgoing: make(chan protocol.Packet, 255),
		Incoming: make(chan protocol.Packet, 255),
		Quit:     make(chan bool),
	}
}

func (cl *Client) Connect(l *Listener) {
	l.ch <- cl.remote
	na := protocol.NewNetworkAdapter(cl.Incoming, cl.Outgoing, cl.Quit, cl.local)
	go cl.readPackets()
	na.Listen()
}

func (cl *Client) ReceivedPackets() (packets []protocol.Packet) {
	cl.packetsMutex.Lock()
	defer cl.packetsMutex.Unlock()

	for _, p := range cl.receivedPackets {
		packets = append(packets, p)
	}

	return
}

func (cl *Client) readPackets() {
	for {
		p := <-cl.Incoming
		cl.packetsMutex.Lock()
		cl.receivedPackets = append(cl.receivedPackets, p)
		cl.packetsMutex.Unlock()
	}
}
