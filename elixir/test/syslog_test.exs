defmodule SyslogTest do
  use ExUnit.Case
  doctest Syslog

  @valid_message ~S(<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8)

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

    assert DateTime.to_iso8601(log.timestamp) == "2003-10-11T22:14:15.003Z"
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
    {log, _, _} = Syslog.parse(~S(<34>1 2003-10-11T22:14:15.003Z - - - - -))
    assert log.hostname == nil
    assert log.app_name == nil
    assert log.proc_id == nil
    assert log.msg_id == nil
    assert log.structure_data == []
    assert log.message == nil
  end

  test "allow escaped characters within structured data values" do
    {log, _, _} =
      Syslog.parse(
        ~S(<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 a="1" b="\"" c="\\" d="\]" e="\"There are \\many things here[1\]\""])
      )

    [sd] = log.structure_data
    assert sd.id == "my@id1"

    assert sd.properties == [
             %Syslog.Property{key: "e", value: ~S("There are \many things here[1]")},
             %Syslog.Property{key: "d", value: "]"},
             %Syslog.Property{key: "c", value: "\\"},
             %Syslog.Property{key: "b", value: "\""},
             %Syslog.Property{key: "a", value: "1"}
           ]
  end

  test "supports multiple sd elements" do
    {log, _, _} =
      Syslog.parse(~S(<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 k="v"][my@id2 c="val"]))

    [sd1, sd2] = log.structure_data
    assert sd1.id == "my@id2"
    assert sd2.id == "my@id1"
  end

  test "example timestamps from the RFC can be parsed" do
    {log, _, _} = Syslog.parse(~S(<34>1 2003-10-11T22:14:15.00003Z - - - - -))
    assert DateTime.to_iso8601(log.timestamp) == "2003-10-11T22:14:15.00003Z"

    {log, _, _} = Syslog.parse(~S(<34>1 1985-04-12T23:20:50.52Z - - - - -))
    assert DateTime.to_iso8601(log.timestamp) == "1985-04-12T23:20:50.52Z"

    {log, _, _} = Syslog.parse(~S(<34>1 1985-04-12T23:20:50.52+00:00 - - - - -))
    assert DateTime.to_iso8601(log.timestamp) == "1985-04-12T23:20:50.52Z"

    {log, _, _} = Syslog.parse(~S(<34>1 1985-04-12T23:20:50.52Z - - - - -))
    assert DateTime.to_iso8601(log.timestamp) == "1985-04-12T23:20:50.52Z"

    {log, _, _} = Syslog.parse(~S(<34>1 1985-04-12T23:20:50.52Z - - - - -))
    assert DateTime.to_iso8601(log.timestamp) == "1985-04-12T23:20:50.52Z"
  end

  test "date can be nil" do
    {log, _, _} = Syslog.parse(~S(<34>1 - - su - - - 'su root' failed for lonvick on /dev/pts/8))
    assert log.timestamp == nil
  end

  test "ignores unparseable timestamps" do
    {log, _, _} = Syslog.parse(~S(<34>1 asdfasdfasdf - su - - - 'su root' failed for lonvick on /dev/pts/8))
    assert log.timestamp == nil
  end
end
