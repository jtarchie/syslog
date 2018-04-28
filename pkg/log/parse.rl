package syslog

import (
  "fmt"
  "time"
  "bytes"
)

%%{
  machine syslog_rfc5424;
  write data;
}%%

func toString(a []byte) string {
  if len(a) == 1 && a[0] == '-' {
    return ""
  }
  return string(a)
}

func atoi(a []byte) int {
  var x int
  for _, c := range a {
    x = x * 10 + int(c - '0')
  }
  return x
}

func atoi2(a []byte) int {
  return int(a[1] - '0') + int(a[0] - '0') * 10
}

func atoi4(a []byte) int {
  return int(a[3] - '0') +
  int(a[2] - '0') * 10 +
  int(a[1] - '0') * 100 +
  int(a[0] - '0') * 1000
}

func Parse(data []byte) (*Log, int, error) {
  var (
    paramName string
    nanosecond int
  )

  log := &Log{}
  var location *time.Location
  var buffer bytes.Buffer

  // set defaults for state machine parsing
  cs, p, pe, eof := 0, 0, len(data), len(data)

  // use to keep track start of value
  mark := 0

  // taken directly from https://tools.ietf.org/html/rfc5424#page-8
  %%{
    action mark      { mark = p }
    action tcp_len   { pe, eof = atoi(data[mark:p]) + (p-mark) + 1, atoi(data[mark:p]) + (p-mark) + 1 }
    action version   { log.version = atoi(data[mark:p]) }
    action priority  { log.priority = atoi(data[mark:p]) }
    action hostname  { log.hostname = toString(data[mark:p]) }
    action appname   { log.appname = toString(data[mark:p]) }
    action procid    { log.procID = toString(data[mark:p]) }
    action msgid     { log.msgID = toString(data[mark:p]) }
    action sdid      {
      log.data = append(log.data, structureElement{
        id: string(data[mark:p]),
        properties: make([]Property, 0, 5),
      })
    }
    action paramname  { paramName = string(data[mark:p]) }
    action escaped    {
      buffer.Write(data[mark:p-2])
      buffer.WriteByte(data[p-1])
      mark = p
    }
    action paramvalue {
      buffer.Write(data[mark:p])
      log.data[len(log.data)-1].properties = append(log.data[len(log.data)-1].properties, Property{paramName,buffer.String()})
      buffer.Reset()
    }

    action timestamp {
      location = time.UTC
      if data[mark+19] == '.' {
        offset := 1

        if data[p-1] != 'Z' {
          offset = 6
          dir := 1
          if data[p-6] == '-' {
            dir = -1
          }

          location = time.FixedZone(
            "",
            dir * (atoi2(data[p-5:p-3]) * 3600 + atoi(data[p-2:p]) * 60),
          )
        }
        nbytes := ( p - offset - 1 ) - ( mark + 19 )
        for i := mark + 20; i < p-offset; i++ {
          nanosecond = nanosecond*10 + int(data[i]-'0')
        }
        for i := 0; i < 9-nbytes; i++ {
          nanosecond *= 10
        }
      }

      log.timestamp = time.Date(
        atoi4(data[mark:mark+4]),
        time.Month(atoi2(data[mark+5:mark+7])),
        atoi2(data[mark+8:mark+10]),
        atoi2(data[mark+11:mark+13]),
        atoi2(data[mark+14:mark+16]),
        atoi2(data[mark+17:mark+19]),
        nanosecond,
        location,
      ).UTC()
    }
    action message { log.message = string(data[mark:p]) }

    include syslog_rfc5424 "syslog.rl";
    write init;
    write exec;
  }%%

  if cs < syslog_rfc5424_first_final {
    return nil, 0, fmt.Errorf("error in msg at pos %d: %s", p, data)
  }

  return log, p, nil
}
