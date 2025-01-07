package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)


const defaultListenAddr = ":5001"


type Config struct {
	ListenAddr string
}

type Server struct {
	Config *Config
	peers map[*peer]bool
	ln net.Listener
	addPeerCh chan *peer
}


func NewServer(cfg *Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config: cfg,
		peers: make(map[*peer]bool),
		addPeerCh: make(chan *peer),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Config.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()
			
	return s.acceptLoop()
	
}


func (s *Server) loop() {
	for {
		select {
		case p := <-s.addPeerCh:
			s.peers[p] = true
			default:
				fmt.Println("foo")
		}
	}
}
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error: %v", err)
			continue
		}
		go s.handleconn(conn)
	}
}

func (s *Server) handleconn (conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer
	peer.readLoop()
}

func main() {
	Server := NewServer(&Config{})
	log.Fatal(Server.Start())
}