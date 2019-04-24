package main

import (
	"fmt"
	"net"

	lua "github.com/yuin/gopher-lua"
)

type udpListener struct {
	port     int
	router   *lua.LFunction
	listener net.PacketConn
	runner   runnerFunc
}

func initializeUDPlistener(
	port int,
	router *lua.LFunction,
	runner runnerFunc,
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
	listener, err := net.ListenPacket("udp", fmt.Sprintf(":%d", u.port))
	if err != nil {
		return fmt.Errorf("could not start server: %s", err)
	}

	u.listener = listener

	go func() {
		buffer := make([]byte, 65535)
	packet:
		n, _, err := u.listener.ReadFrom(buffer)
		if err != nil {
			goto packet
		}
		u.runner(string(buffer[:n]), u.router)
		goto packet
	}()
	return nil
}

func (u *udpListener) Stop() {
	u.listener.Close()
}
