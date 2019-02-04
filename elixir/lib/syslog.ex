defmodule Syslog do
  use Bitwise

  defmodule Property do
    defstruct [:key, :value]
  end

  def parse(msg) do
    {:ok, val, _, _, _, _} = SyslogParser.message(msg)

    IO.inspect(val)

    log = build(%SyslogLog{}, val)
    {log, 0, nil}
  end

  defp build(log, [{:version, version} | p]) do
    log = %{log | version: :erlang.list_to_integer(version)}
    build(log, p)
  end

  defp build(log, [{:prival, [prival]} | p]) do
    log = %{log | severity: prival &&& 7, facility: prival >>> 3, priority: prival}
    build(log, p)
  end

  defp build(log, [{:datetime, datetime} | p]) do
    [year, month, day, hour, minute, second, _, _, _] = datetime

    log = %{
      log
      | timestamp: %DateTime{
          year: year,
          month: month,
          day: day,
          hour: hour,
          minute: minute,
          second: second,
          time_zone: "Etc/UTC",
          zone_abbr: "UTC",
          utc_offset: 0,
          std_offset: 0
        }
    }

    build(log, p)
  end

  defp build(log, [{:hostname, [hostname]} | p]) do
    log = %{log | hostname: hostname}
    build(log, p)
  end

  defp build(log, [{:app_name, [app_name]} | p]) do
    log = %{log | app_name: app_name}
    build(log, p)
  end

  defp build(log, [{:proc_id, [proc_id]} | p]) do
    log = %{log | proc_id: proc_id}
    build(log, p)
  end

  defp build(log, [{:msg_id, [msg_id]} | p]) do
    log = %{log | msg_id: msg_id}
    build(log, p)
  end

  defp build(log, [{:message, [message]} | p]) do
    log = %{log | message: message}
    build(log, p)
  end

  defp build(log, [_ | p]) do
    build(log, p)
  end

  defp build(log, []) do
    log
  end
end
