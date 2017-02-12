package main_test

import (
	"syslog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Syslog", func() {
	msg := []byte("<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")

	It("parses valid messages", func() {
		_, err := syslog.Parser(msg)
		Expect(err).ToNot(HaveOccurred())
	})

	It("sets the version", func() {
		msg, _ := syslog.Parser(msg)
		Expect(msg.Version).To(Equal(uint(1)))
	})

	Measure("message parsing", func(b Benchmarker) {
		runtime := b.Time("10000 messages", func() {
			for i := 0; i < 10000; i++ {
				syslog.Parser(msg)
			}
		})

		Expect(runtime.Seconds()).Should(BeNumerically("<", 1))
	}, 10)
})
