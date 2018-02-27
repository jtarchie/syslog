package syslog

//go:generate ragel -e -G2 -Z parse.rl
// go:generate ragel -Z parse.rl

import (
	"bytes"
	"fmt"
	"time"
)

type Property struct {
	Key   string
	Value string
}

type structureData struct {
	id         string
	properties []Property
}

func (s structureData) ID() string {
	return s.id
}

func (s structureData) Properties() []Property {
	return s.properties
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
	buffer.WriteString(m.timestamp.Format(time.RFC3339Nano))
	buffer.WriteString(" ")
	buffer.WriteString(m.hostname)
	buffer.WriteString(" ")
	buffer.WriteString(m.appname)
	buffer.WriteString(" ")
	buffer.WriteString(m.procID)
	buffer.WriteString(" ")
	buffer.WriteString(m.msgID)
	buffer.WriteString(" ")
	buffer.WriteString("[")
	buffer.WriteString(m.data.id)
	for _, property := range m.data.properties {
		buffer.WriteString(" ")
		buffer.WriteString(property.Key)
		buffer.WriteString(`="`)
		buffer.WriteString(property.Value)
		buffer.WriteString(`"`)
	}
	buffer.WriteString("]")
	buffer.WriteString(" ")
	buffer.WriteString(m.message)
	return buffer.String()
}
