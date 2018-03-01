package transports

import (
	"fmt"
	"log"
	"net"

	"github.com/jtarchie/syslog/pkg/log"
)

type Writer interface {
	Write(*syslog.Log) error
}

type UDPServer struct {
	writer   Writer
	listener net.PacketConn
}

func NewUDPServer(port int, w Writer) (*UDPServer, error) {
	listener, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("could not start server: %s", err)
	}
	return &UDPServer{
		writer:   w,
		listener: listener,
	}, nil
}

func (s *UDPServer) Start() error {
	queue := make(chan []byte, 10000)

	log.Printf("udp: starting server on addr %s", s.listener.LocalAddr().String())
	defer s.listener.Close()

	for i := 1; i <= 10; i++ {
		go func() {
			for {
				select {
				case buffer := <-queue:
					parsed, _, err := syslog.Parse(buffer)
					if err != nil {
						log.Printf("could not parse msg: %s", err)
						continue
					}
					s.writer.Write(parsed)
				}
			}
		}()
	}

	buffer := make([]byte, 1024)
	failed, total := 0, 0
	for {
		n, _, err := s.listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("could not read from UDP: %s", err)
			continue
		}
		total += 1
		select {
		case queue <- buffer[:n]:
		default:
			failed += 1
			if failed%1000 == 0 {
				log.Printf("udp: unable to proccess %d/%d messages with %d in queue", failed, total, len(queue))
			}
		}
	}
}

func (s *UDPServer) Addr() net.Addr {
	return s.listener.LocalAddr()
}
