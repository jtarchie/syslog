package main_test

import (
	"syslog"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Syslog", func() {
	payload := []byte(`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`)

	It("parses valid messages", func() {
		_, err := syslog.Parser(payload)
		Expect(err).ToNot(HaveOccurred())
	})

	It("sets the version", func() {
		msg, _ := syslog.Parser(payload)
		Expect(msg.Version()).To(Equal(1))
	})

	Context("with the priority", func() {
		It("returns the severity", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.Severity()).To(Equal(2))
		})

		It("returns the facility", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.Facility()).To(Equal(4))
		})

		It("sets the priority", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.Priority()).To(Equal(34))
		})
	})

	Context("with the timestamp", func() {
		It("returns a valid date object", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.Timestamp().String()).To(Equal("2003-10-11 22:14:15.003 +0000 UTC"))
		})

		Context("with the example timestamps from the RFC", func() {
			It("can parse them all", func() {
				timestamp := func(t string) *time.Time {
					payload := []byte("<34>1 " + t + " mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")
					msg, err := syslog.Parser(payload)
					Expect(err).ToNot(HaveOccurred())

					return msg.Timestamp()
				}

				Expect(timestamp("1985-04-12T23:20:50.52Z").String()).To(Equal("1985-04-12 23:20:50.52 +0000 UTC"))
				//Expect(timestamp("1985-04-12T19:20:50.52-04:00").String()).To(Equal("1985-04-12 19:20:50.52 -0400 -0400"))
				Expect(timestamp("2003-10-11T22:14:15.003Z").String()).To(Equal("2003-10-11 22:14:15.003 +0000 UTC"))
				//Expect(timestamp("2003-08-24T05:14:15.000003-07:00").String()).To(Equal("2003-08-24 05:14:15.000003 -0700 -0700"))
			})

			It("fails parsing on unsupported formats", func() {
				payload := []byte("<34>1 2003-10-11T22:14:15.003Z07:00 mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")
				_, err := syslog.Parser(payload)
				Expect(err).To(HaveOccurred())
			})
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 - - su - - - 'su root' failed for lonvick on /dev/pts/8")
			msg, _ := syslog.Parser(payload)
			Expect(msg.Timestamp()).To(BeNil())
		})
	})

	Context("with the hostname", func() {
		It("sets the hostname", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.Hostname()).To(BeEquivalentTo("mymachine.example.com"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - su - - - 'su root' failed for lonvick on /dev/pts/8")
			msg, _ := syslog.Parser(payload)
			Expect(msg.Hostname()).To(BeNil())
		})
	})

	Context("with the app name", func() {
		It("sets the app name", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.Appname()).To(BeEquivalentTo("su"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			msg, _ := syslog.Parser(payload)
			Expect(msg.Appname()).To(BeNil())
		})
	})

	Context("with the proc id", func() {
		It("sets the proc id", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.ProcID()).To(BeEquivalentTo("12345"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			msg, _ := syslog.Parser(payload)
			Expect(msg.ProcID()).To(BeNil())
		})
	})

	Context("with the msg id", func() {
		It("sets the msg id", func() {
			msg, _ := syslog.Parser(payload)
			Expect(msg.MsgID()).To(BeEquivalentTo("98765"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			msg, _ := syslog.Parser(payload)
			Expect(msg.MsgID()).To(BeNil())
		})
	})

	Context("with structure data", func() {
		It("sets structure data", func() {
			msg, _ := syslog.Parser(payload)
			data := *msg.StructureData()
			Expect(data.ID()).To(BeEquivalentTo("exampleSDID@32473"))
			Expect(data.Properties()).To(BeEquivalentTo([]syslog.Property{
				{[]byte("iut"), []byte("3")},
				{[]byte("eventSource"), []byte("Application")},
				{[]byte("eventID"), []byte("1011")},
			}))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			msg, _ := syslog.Parser(payload)
			Expect(msg.StructureData()).To(BeNil())
		})
	})

	Measure("message parsing", func(b Benchmarker) {
		runtime := b.Time("10000 messages", func() {
			for i := 0; i < 10000; i++ {
				syslog.Parser(payload)
			}
		})

		Expect(runtime.Seconds()).Should(BeNumerically("<", 1))
	}, 10)
})
