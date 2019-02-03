defmodule Syslog do
  import NimbleParsec
  use Bitwise

  defmodule Log do
    defstruct [
      :version,
      :severity,
      :facility,
      :priority,
      :timestamp
    ]
  end

  date =
    integer(4)
    |> ignore(string("-"))
    |> integer(2)
    |> ignore(string("-"))
    |> integer(2)
  time = integer(2)
    |> ignore(string(":"))
    |> integer(2)
    |> ignore(string(":"))
    |> integer(2)
    |> optional(string("Z"))
  datetime = date |> ignore(string("T")) |> concat(time) |> tag(:datetime)

  prival  = integer(min: 1, max: 3) |> tag(:prival)
  pri     = string("<") |> concat(prival) |> string(">")
  version = ascii_char([?1..?9])
    |> optional(ascii_char([?0..?9]))
    |> optional(ascii_char([?0..?9]))
    |> tag(:version)
  defparsec :message, pri |> concat(version) |> string(" ") |> concat(datetime)

  def parse(msg) do
    {:ok, val, _, _, _, _} = message(msg)

    log = parse1(%Log{}, val)
    {log, 0, nil}
  end

  defp parse1(log, [{:version, version} | p]) do
    log = %{log | version: :erlang.list_to_integer(version) }
    parse1(log, p)
  end

  defp parse1(log, [{:prival, [prival]} | p]) do
    log = %{log |
      severity: prival &&& 7,
      facility: prival >>> 3,
      priority: prival
    }
    parse1(log, p)
  end

  defp parse1(log, [{:datetime, datetime} | p]) do
    [year, month, day, hour, minute, second] = datetime
    log = %{log | timestamp: %DateTime{
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
    }}
    parse1(log, p)
  end

  defp parse1(log, [_ | p]) do
    parse1(log, p)
  end

  defp parse1(log, []) do
    log
  end
end
