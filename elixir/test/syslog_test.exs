defmodule SyslogTest do
  use ExUnit.Case
  doctest Syslog

  @valid_message ~s(<34>1 2003-10-11T22:14:15.003Z)# mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8)

  test "parses a valid message" do
    {_, offset, err} = Syslog.parse(@valid_message)
    assert offset == 0
    assert err == nil
  end

  test "parses the version" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.version == 1
  end

  test "parses priority" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.severity == 2
    assert log.facility == 4
    assert log.priority == 34
  end

  test "parses timestamp" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.timestamp == %DateTime{
        year: 2003,
        month: 10,
        day: 11,
        hour:  22,
        minute: 14,
        second: 15,
        utc_offset: 0,
        time_zone: "Etc/UTC",
        zone_abbr: "UTC",
        std_offset: 0,
      }
  end
end
