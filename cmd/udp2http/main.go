package main

import (
	"log"

	"github.com/jtarchie/syslog"
	"github.com/jtarchie/syslog/web"
)

func main() {
	log.Println("starting servers")
	writer := web.NewServer(8081)
	go func() {
		log.Fatalf("Could not start writer: %s", writer.Start())
	}()

	server, err := syslog.NewUDPServer(8088, writer)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

	server.Start()
}
