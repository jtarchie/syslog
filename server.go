package syslog

import (
	"log"
	"net"
)

type Writer interface {
	Write(*Log) error
}

type Server struct {
	writer   Writer
	listener net.PacketConn
}

func NewServer(w Writer) *Server {
	listener, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatalf("could not start server: %s", err)
	}
	return &Server{
		writer:   w,
		listener: listener,
	}
}

func (s *Server) Start() error {
	defer s.listener.Close()

	buffer := make([]byte, 1024)
	for {
		_, _, err := s.listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("could not read from UDP: %s", err)
			continue
		}

		parsed, err := Parse(buffer)
		if err != nil {
			log.Printf("could not parse msg: %s", err)
			continue
		}
		s.writer.Write(parsed)
	}
	return nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) Addr() net.Addr {
	return s.listener.LocalAddr()
}

func (s *Server) handle(conn net.Conn) {
}
