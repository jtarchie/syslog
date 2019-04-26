package main

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"log"
	"net"
)

type udpListener struct {
	port     int
	router   *lua.LFunction
	listener net.PacketConn
	runner   writerFunc
}

func initializeUDPlistener(
	port int,
	router *lua.LFunction,
	runner writerFunc,
) (*udpListener, error) {
	if port <= 0 {
		return nil, fmt.Errorf("port number (%d) must be greater than 0 for udp server", port)
	}
	return &udpListener{
		port:   port,
		router: router,
		runner: runner,
	}, nil
}

func (u *udpListener) Start() error {
	listener, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%d", u.port))
	if err != nil {
		return fmt.Errorf("could not start server: %s", err)
	}
	log.Printf("started udp server on port %s", listener.LocalAddr().String())

	u.listener = listener

	go func() {
		for {
			log.Println("waiting for packet")
			buffer := make([]byte, 65535)
			n, _, err := u.listener.ReadFrom(buffer)
			if err != nil {
				log.Println("could not read packet: %s", err)
				continue
			}

			u.runner(buffer[:n], u.router)
		}
	}()
	return nil
}

func (u *udpListener) Stop() {
	log.Printf("stopping udp server on port %d", u.port)
	u.listener.Close()
}
