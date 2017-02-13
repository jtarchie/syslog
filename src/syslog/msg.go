package syslog

//go:generate ragel -Z syslog.rl

import "time"

type structureData struct {
	id         string
	properties map[string]string
}

func (s *structureData) ID() string {
	return s.id
}

func (s *structureData) Properties() map[string]string {
	return s.properties
}

type message struct {
	version   int
	priority  int
	timestamp *time.Time
	hostname  *string
	appname   *string
	procID    *string
	msgID     *string
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

func (m *message) Hostname() *string {
	return m.hostname
}

func (m *message) Appname() *string {
	return m.appname
}

func (m *message) ProcID() *string {
	return m.procID
}

func (m *message) MsgID() *string {
	return m.msgID
}

func (m *message) StructureData() *structureData {
	return m.data
}
