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

	It("renders back to its original format", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.String()).To(Equal(validMessage))
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

	It("returns a valid date object", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.Timestamp().String()).To(Equal("2003-10-11 22:14:15.003 +0000 UTC"))
	})

	It("sets the hostname", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.Hostname()).To(BeEquivalentTo("mymachine.example.com"))
	})

	It("sets the app name", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.Appname()).To(BeEquivalentTo("su"))
	})

	It("sets the proc id", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.ProcID()).To(BeEquivalentTo("12345"))
	})

	It("sets the log id", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.MsgID()).To(BeEquivalentTo("98765"))
	})

	It("sets structure data", func() {
		log, _, _ := syslog.Parse(payload)
		data := log.StructureData()[0]
		Expect(data.ID()).To(BeEquivalentTo("exampleSDID@32473"))
		Expect(data.Properties()).To(BeEquivalentTo([]syslog.Property{
			{"iut", "3"},
			{"eventSource", "Application"},
			{"eventID", "1011"},
		}))
	})

	It("sets the message", func() {
		log, _, _ := syslog.Parse(payload)
		Expect(log.Message()).To(BeEquivalentTo("'su root' failed for lonvick on /dev/pts/8"))
	})

}

var _ = Describe("with a standard payload", func() {
	It("returns the correct offset", func() {
		_, offset, _ := syslog.Parse([]byte(validMessage))
		Expect(offset).To(Equal(len(validMessage)))
	})

	ParseLogMessageTests(validMessage)

	Context("with the hostname", func() {
		It("is empty when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - su - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Hostname()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

	Context("with the app name", func() {
		It("is empty when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Appname()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

	Context("with the proc id", func() {
		It("is empty when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.ProcID()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

	Context("with the log id", func() {
		It("is empty when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.MsgID()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

	Context("with structure data", func() {
		It("is empty when '-'", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.StructureData()).To(HaveLen(0))
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

	Context("with a message", func() {
		It("sets empty for no message", func() {
			payload := []byte("<34>1 2003-10-11T22:14:15.003Z - - - - -")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Message()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

	Context("with multiple structure data elements", func() {
		It("parses them all", func() {
			payload := []byte(`<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 k="v"][my@id2 c="val"]`)
			log, _, err := syslog.Parse(payload)
			Expect(err).ToNot(HaveOccurred())
			data := log.StructureData()[0]
			Expect(data.ID()).To(BeEquivalentTo("my@id1"))
			Expect(data.Properties()).To(BeEquivalentTo([]syslog.Property{
				{"k", "v"},
			}))
			data = log.StructureData()[1]
			Expect(data.ID()).To(BeEquivalentTo("my@id2"))
			Expect(data.Properties()).To(BeEquivalentTo([]syslog.Property{
				{"c", "val"},
			}))
		})
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

		It("is sets empty when '-'", func() {
			payload := []byte("<34>1 - - su - - - 'su root' failed for lonvick on /dev/pts/8")
			log, _, _ := syslog.Parse(payload)
			Expect(log.Timestamp().IsZero()).To(BeTrue())
			Expect(log.String()).To(BeEquivalentTo(payload))
		})
	})

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

	It("works with these messages", func() {
		_, _, err := syslog.Parse([]byte(`<4>1 2018-02-02T18:30:55.968Z worker-25.example.com flood 64974 937 - This is message 937 from worker`))
		Expect(err).ToNot(HaveOccurred())
	})
})
