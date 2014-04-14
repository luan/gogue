package fakes

import (
	"net"

	"github.com/luan/gogue/protocol"
)

type Client struct {
	remote   *Conn
	local    *Conn
	Outgoing chan protocol.Packet
	Incoming chan protocol.Packet
	Quit     chan bool
}

func NewClient() *Client {
	serverBuffer, clientBuffer := net.Pipe()
	return &Client{
		remote:   NewConn(serverBuffer),
		local:    NewConn(clientBuffer),
		Outgoing: make(chan protocol.Packet),
		Incoming: make(chan protocol.Packet),
		Quit:     make(chan bool),
	}
}

func (cl *Client) Connect(l *Listener) {
	l.ch <- cl.remote
	na := protocol.NewNetworkAdapter(cl.Incoming, cl.Outgoing, cl.Quit, cl.local)
	na.Listen()
}
