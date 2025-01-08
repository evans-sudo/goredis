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
	quitCh chan struct{}
	msgCh chan []byte
}


func NewServer(cfg *Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config: cfg,
		peers: make(map[*peer]bool),
		addPeerCh: make(chan *peer),
		quitCh: make(chan struct{}),
		msgCh: make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Config.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	//slog.Info("Listening on", s.Config.ListenAddr)
			
	return s.acceptLoop()	
}

func (s *Server) handleRawmsg(rawMsg []byte) error {
	fmt.Println(string(rawMsg))
	return nil
}


func (s *Server) loop() {
	for {
		select {
		case Rawmsg := <-s.msgCh:
			if err := s.handleRawmsg(Rawmsg); err != nil {
			slog.Info("raw message error", "err", err)
			}
			fmt.Println("msg received", (Rawmsg))
			
		case <-s.quitCh:
			return
		case p := <-s.addPeerCh:
			s.peers[p] = true
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
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("New peer","remoteaddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteaddr", conn.RemoteAddr())
	}
}

func main() {
	Server := NewServer(&Config{})
	log.Fatal(Server.Start())
}