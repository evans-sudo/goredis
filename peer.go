package main

import (
	"net"
)

type peer struct {
	conn net.Conn
	msgch chan []byte
}


func NewPeer(conn net.Conn, msgch chan []byte) *peer {
	return &peer{
		conn: conn,
		msgch: msgch,
	}
}

func (p *peer) readLoop() error {
	buf := make([]byte, 1024)
	for {
		// read from conn
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgch <- msgBuf
	}
}