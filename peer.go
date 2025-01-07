package main

import "net"

type peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *peer {
	return &peer{
		conn: conn,
	}
}

func (p *peer) readLoop() {
	for {
		// read from conn
	}
}