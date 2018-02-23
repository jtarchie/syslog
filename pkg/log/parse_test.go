package syslog_test

import (
	"fmt"
	"time"

	"github.com/jtarchie/syslog/pkg/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const validMessage = `<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`

func ParseLogMessageTests(message string) {
	var (
		payload []byte
		log     *syslog.Log
	)

	BeforeEach(func() {
		payload = []byte(message)
		log = &syslog.Log{}
	})

	It("parses valid messages", func() {
		_, offset, err := syslog.Parse(payload)
		Expect(offset).ToNot(Equal(0))
		Expect(err).ToNot(HaveOccurred())
	})

	It("sets the version", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.Version()).To(Equal(1))
	})

	Context("with the priority", func() {
		It("returns the severity", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Severity()).To(Equal(2))
		})

		It("returns the facility", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Facility()).To(Equal(4))
		})

		It("sets the priority", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Priority()).To(Equal(34))
		})
	})

	Context("with the timestamp", func() {
		It("returns a valid date object", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Timestamp().String()).To(Equal("2003-10-11 22:14:15.003 +0000 UTC"))
		})

		Context("with the example timestamps from the RFC", func() {
			It("can parse them all", func() {
				timestamp := func(t string) time.Time {
					payload := []byte("<34>1 " + t + " mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")
					log, _, err := syslog.Parse(payload)
					Expect(err).ToNot(HaveOccurred())

					return log.Timestamp()
				}

				Expect(timestamp("2003-10-11T22:14:15.00003Z").String()).To(Equal("2003-10-11 22:14:15.00003 +0000 UTC"))
				Expect(timestamp("1985-04-12T23:20:50.52Z").String()).To(Equal("1985-04-12 23:20:50.52 +0000 UTC"))
				Expect(timestamp("1985-04-12T23:20:50.52+00:00").String()).To(Equal("1985-04-12 23:20:50.52 +0000 UTC"))
				Expect(timestamp("1985-04-12T23:20:50.52+02:00").String()).To(Equal("1985-04-12 21:20:50.52 +0000 UTC"))
				Expect(timestamp("1985-04-12T18:20:50.52-02:00").String()).To(Equal("1985-04-12 20:20:50.52 +0000 UTC"))
			})

			It("fails parsing on unsupported formats", func() {
				payload := []byte("<34>1 2003-10-11T22:14:15.003Z07:00 mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8")
				_, n, err := syslog.Parse(payload)
				Expect(n).To(Equal(0))
				Expect(err).To(HaveOccurred())
			})
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 - - su - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Timestamp().IsZero()).To(BeTrue())
		})
	})

	Context("with the hostname", func() {
		It("sets the hostname", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Hostname()).To(BeEquivalentTo("mymachine.example.com"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - su - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Hostname()).To(BeEmpty())
		})
	})

	Context("with the app name", func() {
		It("sets the app name", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Appname()).To(BeEquivalentTo("su"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Appname()).To(BeEmpty())
		})
	})

	Context("with the proc id", func() {
		It("sets the proc id", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.ProcID()).To(BeEquivalentTo("12345"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.ProcID()).To(BeEmpty())
		})
	})

	Context("with the log id", func() {
		It("sets the log id", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.MsgID()).To(BeEquivalentTo("98765"))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.MsgID()).To(BeEmpty())
		})
	})

	Context("with structure data", func() {
		It("sets structure data", func() {
			log, _, _ := syslog.Parse(payload)
			data := log.StructureData()
			Expect(data.ID()).To(BeEquivalentTo("exampleSDID@32473"))
			Expect(data.Properties()).To(BeEquivalentTo([]syslog.Property{
				{"iut", "3"},
				{"eventSource", "Application"},
				{"eventID", "1011"},
			}))
		})

		It("is nil when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.StructureData().ID()).To(BeEmpty())
			Expect(log.StructureData().Properties()).To(BeNil())
		})
	})

	Context("with a message", func() {
		It("sets the message", func() {
			log, _, _ := syslog.Parse(payload)
			Expect(log.Message()).To(BeEquivalentTo("'su root' failed for lonvick on /dev/pts/8"))
		})

		It("sets nil for no message", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - -")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Message()).To(BeEmpty())
		})
	})
}

var _ = Describe("with a standard payload", func() {
	It("returns the correct offset", func() {
		_, offset, _ := syslog.Parse(payload)
		Expect(offset).To(Equal(len(validMessage)))
	})

	ParseLogMessageTests(validMessage)
})

var _ = Describe("with a TCP payload", func() {
	Context("with a single message", func() {
		payload := fmt.Sprintf("%d %s", len(validMessage), validMessage)

		It("returns the correct offset", func() {
			_, offset, _ := syslog.Parse([]byte(payload))
			Expect(offset).To(Equal(len(validMessage) + 4))
		})

		ParseLogMessageTests(payload)
	})

	Context("with a multiple payload", func() {
		payload := fmt.Sprintf("%d %s%d %s", len(validMessage), validMessage, len(validMessage), validMessage)

		It("returns the correct offset", func() {
			_, offset, _ := syslog.Parse([]byte(payload))
			Expect(offset).To(Equal(len(validMessage) + 4))
		})
		ParseLogMessageTests(payload)
	})
})
