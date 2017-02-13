package syslog

import (
  "fmt"
  "strconv"
  "time"
)

%%{
  machine syslog_rfc5424;
  write data;
}%%

var timestampFormats = []string{
  "2006-01-02T15:04:05.999999999-07:00",
  "2006-01-02T15:04:05-07:00",
  "2006-01-02T15:04:05.999999999Z",
  "2006-01-02T15:04:05Z",
}

func stringPtr(d []byte) *string {
  s := string(d)
  if s == "-" {
    return nil
  }
  return &s
}

func Parser(data []byte) (*message, error) {
  msg := &message{}

  // set defaults for state machine parsing
  cs, p, pe := 0, 0, len(data)

  // use to keep track start of value
  mark, paramName := 0, ""

  // taken directly from https://tools.ietf.org/html/rfc5424#page-8
  %%{
    action mark       { mark = p }
    action version    { msg.version, _ = strconv.Atoi(string(data[mark:p])) }
    action priority   { msg.priority, _ = strconv.Atoi(string(data[mark:p])) }
    action timestamp  {
      for _, format := range timestampFormats {
        t, err := time.Parse(format, string(data[mark:p]))
        if err == nil {
          msg.timestamp = &t
          break
        }
      }
    }
    action hostname   { msg.hostname = stringPtr(data[mark:p]) }
    action appname    { msg.appname = stringPtr(data[mark:p]) }
    action procid     { msg.procID = stringPtr(data[mark:p]) }
    action msgid      { msg.msgID = stringPtr(data[mark:p]) }
    action sdid       {
      msg.data = &structureData{
        id: string(data[mark:p]),
        properties: make(map[string]string),
      }
    }
    action paramname  { paramName = string(data[mark:p]) }
    action paramvalue { msg.data.properties[paramName] = string(data[mark:p]) }

    nil           = "-";
    nonzero_digit = "1".."9";
    printusascii  = "!".."~";
    sp            = " ";
    octet         = 0..255;

    utf8_string = octet*;
    bom         = 0xEF 0xBB 0xBF;
    msg_utf8    = bom utf8_string;
    msg_any     = octet*;
    msg         = msg_any | msg_utf8;

    sd_name         = printusascii{1,32} -- ("=" | sp | "]" | 34 | '"');
    param_value     = utf8_string >mark %paramvalue;
    param_name      = sd_name >mark %paramname;
    sd_id           = sd_name >mark %sdid;
    sd_param        = param_name "=" '"' param_value '"';
    sd_element      = "[" sd_id ( sp sd_param )* "]";
    structured_data = nil | sd_element{,1};

    time_hour      = digit{,2};
    time_minute    = digit{,2};
    time_second    = digit{,2};
    time_secfrac   = "." digit{1,6};
    time_numoffset = ("+" | "-") time_hour ":" time_minute;
    time_offset    = "Z" | time_numoffset;
    partial_time   = time_hour ":" time_minute ":" time_second time_secfrac?;
    full_time      = partial_time time_offset;
    date_mday      = digit{,2};
    date_month     = digit{,2};
    date_fullyear  = digit{,4};
    full_date      = date_fullyear "-" date_month "-" date_mday;
    timestamp      = nil | (full_date "T" full_time) >mark %timestamp;

    msg_id   = (nil | printusascii{1,32}) >mark %msgid;
    proc_id  = (nil | printusascii{1,128}) >mark %procid;
    app_name = (nil | printusascii{1,48}) >mark %appname;

    hostname = (nil | printusascii{1,255}) >mark %hostname;
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
