package syslog

//go:generate ragel -e -G2 -Z parse.rl
// go:generate ragel -Z parse.rl

import "time"

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
