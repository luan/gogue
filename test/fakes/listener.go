package fakes

import "net"

type Listener struct {
	ch chan net.Conn
}

func NewListener() *Listener {
	return &Listener{
		ch: make(chan net.Conn),
	}
}

func (l *Listener) Accept() (c net.Conn, err error) {
	c = <-l.ch
	return
}

func (l *Listener) Addr() (a net.Addr) {
	return
}

func (l *Listener) Close() (err error) {
	return
}
