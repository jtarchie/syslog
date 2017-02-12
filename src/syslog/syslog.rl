package syslog

import "fmt"

type message struct{
  Version uint64
}

%%{
  machine syslog_rfc5424;
  write data;
}%%

func Parser(data []byte) (*message, error) {
  msg := &message{}

  // set defaults for state machine parsing
  cs, p, pe := 0, 0, len(data)

  // taken directly from https://tools.ietf.org/html/rfc5424#page-8
  %%{
    nil           = "-";
    nonzero_digit = 49..57;
    printusascii  = 33..126;
    sp            = 32;
    octet         = 0..255;

    utf8_string = octet*;
    bom         = 0xEF 0xBB 0xBF;
    msg_utf8    = bom utf8_string;
    msg_any     = octet*;
    msg         = msg_any | msg_utf8;

    sd_name         = printusascii {1,32} -- ("=" | sp | "]" | 34 | '"');
    param_value     = utf8_string;
    param_name      = sd_name;
    sd_id           = sd_name;
    sd_element      = param_name "=" 34 param_value 34;
    structured_data = nil | sd_element+;

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
    timestamp      = nil | full_date "T" full_time;

    msg_id   = nil | printusascii{1,32};
    proc_id  = nil | printusascii{1,128};
    app_name = nil | printusascii{1,48};

    hostname = nil | printusascii{1,255};
    version  = nonzero_digit digit{0,2};
    prival   = digit{1,3};
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
