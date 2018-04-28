package syslog

//go:generate ragel -e -G2 -Z parse.rl

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

var paramValueReplacer = strings.NewReplacer(`"`, `\"`, `]`, `\]`, `\`, `\\`)

type Property struct {
	Key   string
	Value string
}

type structureElement struct {
	id         string
	properties []Property
}

func (s structureElement) ID() string {
	return s.id
}

func (s structureElement) Properties() []Property {
	return s.properties
}

func (s structureElement) String() string {
	if s.id == "" {
		return "-"
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(s.id)
	for _, property := range s.properties {
		buffer.WriteString(" ")
		buffer.WriteString(property.Key)
		buffer.WriteString(`="`)
		buffer.WriteString(paramValueReplacer.Replace(property.Value))
		buffer.WriteString(`"`)
	}
	buffer.WriteString("]")
	return buffer.String()
}

type structureData []structureElement

func (s structureData) String() string {
	if len(s) == 0 {
		return "-"
	}

	var buffer bytes.Buffer
	for _, element := range s {
		buffer.WriteString(element.String())
	}
	return buffer.String()
}

type Log struct {
	version   int
	priority  int
	timestamp time.Time
	hostname  string
	appname   string
	procID    string
	msgID     string
	data      structureData
	message   string
}

func (m *Log) Version() int {
	return m.version
}

func (m *Log) Severity() int {
	return m.priority & 7
}

func (m *Log) Facility() int {
	return m.priority >> 3
}

func (m *Log) Priority() int {
	return m.priority
}

func (m *Log) Timestamp() time.Time {
	return m.timestamp
}

func (m *Log) Hostname() string {
	return m.hostname
}

func (m *Log) Appname() string {
	return m.appname
}

func (m *Log) ProcID() string {
	return m.procID
}

func (m *Log) MsgID() string {
	return m.msgID
}

func (m *Log) StructureData() structureData {
	return m.data
}

func (m *Log) Message() string {
	return m.message
}

func (m *Log) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("<")
	buffer.WriteString(fmt.Sprintf("%d", m.priority))
	buffer.WriteString(">")
	buffer.WriteString(fmt.Sprintf("%d", m.version))
	buffer.WriteString(" ")
	if m.timestamp.IsZero() {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.timestamp.Format(time.RFC3339Nano))
	}
	buffer.WriteString(" ")
	if m.hostname == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.hostname)
	}
	buffer.WriteString(" ")
	if m.appname == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.appname)
	}
	buffer.WriteString(" ")
	if m.procID == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.procID)
	}
	buffer.WriteString(" ")
	if m.msgID == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.msgID)
	}
	buffer.WriteString(" ")
	buffer.WriteString(m.data.String())
	if m.message != "" {
		buffer.WriteString(" ")
		buffer.WriteString(m.message)
	}
	return buffer.String()
}
