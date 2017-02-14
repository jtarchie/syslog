package syslog

import (
  "fmt"
  "bytes"
  "time"
)

%%{
  machine syslog_rfc5424;
  write data;
}%%

var nilValue = []byte("-")

func bytesRef(a []byte) []byte {
  if bytes.Compare(a, nilValue) == 0 {
    return nil
  }
  return a
}

func atoi(a []byte) int {
  var x int
  for _, c := range a {
    x = x * 10 + int(c - '0')
  }
  return x
}

func power(value, times int) int {
  for i := 0; i < times; i++ {
    value *= 10
  }
  return value
}

type timestamp struct{
  year, month, day, hour, minute, second, nsec int
}

func Parser(data []byte) (*message, error) {
  var (
    paramName []byte
    t time.Time
  )
  msg, timestamp := &message{}, timestamp{}

  // set defaults for state machine parsing
  cs, p, pe := 0, 0, len(data)

  // use to keep track start of value
  mark := 0

  // taken directly from https://tools.ietf.org/html/rfc5424#page-8
  %%{
    action mark      { mark = p }
    action version   { msg.version = atoi(data[mark:p]) }
    action priority  { msg.priority = atoi(data[mark:p]) }
    action hostname  { msg.hostname = bytesRef(data[mark:p]) }
    action appname   { msg.appname = bytesRef(data[mark:p]) }
    action procid    { msg.procID = bytesRef(data[mark:p]) }
    action msgid     { msg.msgID = bytesRef(data[mark:p]) }
    action sdid      {
      msg.data = &structureData{
        id: data[mark:p],
        properties: []Property{},
      }
    }
    action paramname  { paramName = data[mark:p] }
    action paramvalue { msg.data.properties = append(msg.data.properties, Property{paramName,data[mark:p]}) }

    action year    { timestamp.year   = atoi(data[mark:p]) }
    action month   { timestamp.month  = atoi(data[mark:p]) }
    action mday    { timestamp.day    = atoi(data[mark:p]) }
    action hour    { timestamp.hour   = atoi(data[mark:p]) }
    action minute  { timestamp.minute = atoi(data[mark:p]) }
    action second  { timestamp.second = atoi(data[mark:p]) }
    action secfrac { timestamp.nsec   = power(atoi(data[mark:p]), 9-(p-mark)) }

    action timestamp {
      t = time.Date(
        timestamp.year,
        time.Month(timestamp.month),
        timestamp.day,
        timestamp.hour,
        timestamp.minute,
        timestamp.second,
        timestamp.nsec,
        time.UTC,
      )
      msg.timestamp = &t
    }

    nil           = "-";
    nonzero_digit = "1".."9";
    printusascii  = "!".."~";
    sp            = " ";

    utf8_string = any*;
    bom         = 0xEF 0xBB 0xBF;
    msg_utf8    = bom utf8_string;
    msg_any     = any*;
    msg         = msg_any | msg_utf8;

    sd_name         = printusascii{1,32} -- ("=" | sp | "]" | 34 | '"');
    param_value     = utf8_string >mark %paramvalue;
    param_name      = sd_name >mark %paramname;
    sd_id           = sd_name >mark %sdid;
    sd_param        = param_name '="' param_value :>> '"';
    sd_element      = "[" sd_id ( sp sd_param )* "]";
    structured_data = nil | sd_element{,1};

    time_hour      = digit{,2};
    time_minute    = digit{,2};
    time_second    = digit{,2};
    time_secfrac   = "." (digit{1,6} >mark %secfrac);
    #time_numoffset = ("+" | "-") time_hour ":" time_minute;
    time_offset    = "Z"; #| time_numoffset;
    partial_time   = (time_hour >mark %hour)
      ":" (time_minute >mark %minute)
      ":" (time_second >mark %second)
      time_secfrac?;
    full_time      = partial_time time_offset;
    date_mday      = digit{,2} >mark %mday;
    date_month     = digit{,2} >mark %month;
    date_fullyear  = digit{,4} >mark %year;
    full_date      = date_fullyear "-" date_month "-" date_mday;
    timestamp      = nil | (full_date "T" full_time %timestamp);

    msg_id   = nil | printusascii{1,32} >mark %msgid;
    proc_id  = nil | printusascii{1,128} >mark %procid;
    app_name = nil | printusascii{1,48} >mark %appname;

    hostname = nil | printusascii{1,255} >mark %hostname;
    version  = (nonzero_digit digit{0,2}) >mark %version;
    prival   = digit{1,3} >mark %priority;
    pri      = "<" prival ">";
    header   = pri version sp timestamp sp hostname sp app_name sp proc_id sp msg_id;

    syslog_msg = header sp structured_data (sp msg)?;
    main := syslog_msg;

    write init;
    write exec;

  }%%

  if cs < syslog_rfc5424_first_final {
    return nil, fmt.Errorf("error in msg at pos %d: %s", p, data)
  }

  return msg, nil
}
