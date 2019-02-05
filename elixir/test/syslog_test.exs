defmodule SyslogTest do
  use ExUnit.Case
  doctest Syslog

  @valid_message ~s(<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8)

  test "parses a valid message" do
    {_, offset, err} = Syslog.parse(@valid_message)
    assert offset == 0
    assert err == nil
  end

  test "sets the version" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.version == 1
  end

  test "sets priority (and decoded fields)" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.severity == 2
    assert log.facility == 4
    assert log.priority == 34
  end

  test "returns a valid date object" do
    {log, _, _} = Syslog.parse(@valid_message)

    assert log.timestamp == %DateTime{
             year: 2003,
             month: 10,
             day: 11,
             hour: 22,
             minute: 14,
             second: 15,
             utc_offset: 0,
             time_zone: "Etc/UTC",
             zone_abbr: "UTC",
             std_offset: 0
           }
  end

  test "sets the hostname" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.hostname == "mymachine.example.com"
  end

  test "sets the app name" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.app_name == "su"
  end

  test "sets the proc id" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.proc_id == "12345"
  end

  test "sets the log id" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.msg_id == "98765"
  end

  test "sets structure data" do
    {log, _, _} = Syslog.parse(@valid_message)
    [sd] = log.structure_data
    assert sd.id == "exampleSDID@32473"

    assert sd.properties == [
             %Syslog.Property{key: "eventID", value: "1011"},
             %Syslog.Property{key: "eventSource", value: "Application"},
             %Syslog.Property{key: "iut", value: "3"}
           ]
  end

  test "sets the message" do
    {log, _, _} = Syslog.parse(@valid_message)
    assert log.message == "'su root' failed for lonvick on /dev/pts/8"
  end

  test "sets values to null when set to '-'" do
    {log, _, _} = Syslog.parse(~s(<34>1 2003-10-11T22:14:15.003Z - - - - -))
    assert log.hostname == nil
    assert log.app_name == nil
    assert log.proc_id == nil
    assert log.msg_id == nil
    assert log.structure_data == []
    assert log.message == nil
  end
end
