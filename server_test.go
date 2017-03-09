package syslog_test

import (
	"net"

	"github.com/jtarchie/syslog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func SendMessage(s *syslog.Server) {
	payload := []byte(`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`)

	conn, err := net.Dial("udp", s.Addr().String())
	Expect(err).ToNot(HaveOccurred())
	defer conn.Close()

	_, err = conn.Write(payload)
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Server", func() {
	It("accepts datagrams via UDP", func() {
		writer := &SpyWriter{}
		server := syslog.NewServer(writer)
		go server.Start()
		defer server.Close()

		SendMessage(server)

		Eventually(func() int {
			return len(writer.Logs)
		}).Should(Equal(1))
		Expect(writer.Logs[0].Message()).To(BeEquivalentTo(`'su root' failed for lonvick on /dev/pts/8`))
	})
})

type SpyWriter struct {
	Logs []*syslog.Log
}

func (s *SpyWriter) Write(l *syslog.Log) error {
	s.Logs = append(s.Logs, l)
	return nil
}
