package syslog_test

import (
	"net"
	"testing"

	"github.com/jtarchie/syslog"
)

var payload = []byte(`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`)
var err error

func BenchmarkSendMessage(b *testing.B) {
	writer := &SpyWriter{}
	server := syslog.NewServer(writer)
	go server.Start()

	conn, _ := net.Dial("udp", server.Addr().String())
	defer conn.Close()

	conn1, _ := net.Dial("udp", server.Addr().String())
	defer conn1.Close()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		conn.Write(payload)
		conn1.Write(payload)
	}
}

func BenchmarkParsing(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err = syslog.Parse(payload)
	}
}
