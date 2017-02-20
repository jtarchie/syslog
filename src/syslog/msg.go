package syslog

//go:generate ragel -e -G2 -Z syslog.rl

import "time"

type Property struct {
	Key   []byte
	Value []byte
}

type structureData struct {
	id         []byte
	properties []Property
}

func (s *structureData) ID() []byte {
	return s.id
}

func (s *structureData) Properties() []Property {
	return s.properties
}

type message struct {
	version   int
	priority  int
	timestamp *time.Time
	hostname  []byte
	appname   []byte
	procID    []byte
	msgID     []byte
	data      *structureData
}

func (m *message) Version() int {
	return m.version
}

func (m *message) Severity() int {
	return m.priority & 7
}

func (m *message) Facility() int {
	return m.priority >> 3
}

func (m *message) Priority() int {
	return m.priority
}

func (m *message) Timestamp() *time.Time {
	return m.timestamp
}

func (m *message) Hostname() []byte {
	return m.hostname
}

func (m *message) Appname() []byte {
	return m.appname
}

func (m *message) ProcID() []byte {
	return m.procID
}

func (m *message) MsgID() []byte {
	return m.msgID
}

func (m *message) StructureData() *structureData {
	return m.data
}
