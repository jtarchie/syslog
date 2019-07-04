package syslog

import (
  "time"
  "errors"
)

%%{
  machine syslog_rfc5424;
  write data;
}%%

var parseError = errors.New("could not parse message")

const empty = ""

func toString(a []byte) string {
  if len(a) == 1 && a[0] == '-' {
    return empty
  }
  return string(a)
}

func atoi(a []byte) int {
  var x, i int
loop:
  x = x * 10 + int(a[i] - '0')
  i++
  if i < len(a) {
    goto loop // avoid for loop so this function can be inlined
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
    ts int
  )

  log := Log{}
  var location *time.Location
  var buffer []byte

  // set defaults for state machine parsing
  cs, p, pe, eof := 0, 0, len(data), len(data)

  // taken directly from https://tools.ietf.org/html/rfc5424#page-8
  %%{
    action tcp_len   { pe, eof = atoi(data[ts:p]) + (p-ts) + 1, atoi(data[ts:p]) + (p-ts) + 1 }
    action version   { log.version = atoi(data[ts:p]) }
    action priority  { log.priority = atoi(data[ts:p]) }
    action hostname  { log.hostname = toString(data[ts:p]) }
    action appname   { log.appname = toString(data[ts:p]) }
    action procid    { log.procID = toString(data[ts:p]) }
    action msgid     { log.msgID = toString(data[ts:p]) }
    action sdid      {
      log.data = append(log.data, structureElement{
        id: string(data[ts:p]),
        properties: make([]Property, 0, 5),
      })
    }
    action paramname  { paramName = string(data[ts:p]) }
    action escaped    {
      buffer = append(buffer, data[ts:p-2]...)
      buffer = append(buffer, data[p-1])
      ts = p
    }
    action paramvalue {
      buffer = append(buffer, data[ts:p]...)
      log.data[len(log.data)-1].properties = append(log.data[len(log.data)-1].properties, Property{paramName,string(buffer)})
      buffer = buffer[:0]
    }

    action timestamp {
      location = time.UTC
      if data[ts+19] == '.' {
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
        nbytes := ( p - offset - 1 ) - ( ts + 19 )
        i := ts + 20
        first:
          if i < p-offset {
            nanosecond = nanosecond*10 + int(data[i]-'0')
            i++
            goto first
          }
        i = 0
        second:
          if i < 9-nbytes {
            nanosecond *= 10
            i++
            goto second
          }
      }

      log.timestamp = time.Date(
        atoi4(data[ts:ts+4]),
        time.Month(atoi2(data[ts+5:ts+7])),
        atoi2(data[ts+8:ts+10]),
        atoi2(data[ts+11:ts+13]),
        atoi2(data[ts+14:ts+16]),
        atoi2(data[ts+17:ts+19]),
        nanosecond,
        location,
      ).UTC()
    }
    action message { log.message = string(data[ts:p]) }

    include syslog_rfc5424 "syslog.rl";
    write init;
    write exec;
  }%%

  if cs < syslog_rfc5424_first_final {
    return nil, 0, parseError
  }

  return &log, p, nil
}
