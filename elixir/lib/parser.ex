defmodule SyslogParser do
  import NimbleParsec

  date =
    integer(4)
    |> ignore(string("-"))
    |> integer(2)
    |> ignore(string("-"))
    |> integer(2)

  time =
    integer(2)
    |> ignore(string(":"))
    |> integer(2)
    |> ignore(string(":"))
    |> integer(2)
    |> choice([
      optional(string(".") |> ascii_string([?0..?9], max: 6) |> optional(string("Z"))),
      optional(string("Z"))
    ])

  datetime = date |> ignore(string("T")) |> concat(time) |> tag(:datetime)

  prival = integer(min: 1, max: 3) |> tag(:prival)
  pri = ignore(string("<")) |> concat(prival) |> ignore(string(">"))

  version =
    ascii_char([?1..?9])
    |> optional(ascii_char([?0..?9]))
    |> optional(ascii_char([?0..?9]))
    |> tag(:version)

  hostname =
    choice([
      ignore(string("-")),
      ascii_string([?!..?~], max: 255)
    ])
    |> tag(:hostname)

  app_name =
    choice([
      ignore(string("-")),
      ascii_string([?!..?~], max: 48)
    ])
    |> tag(:app_name)

  proc_id =
    choice([
      ignore(string("-")),
      ascii_string([?!..?~], max: 128)
    ])
    |> tag(:proc_id)

  msg_id =
    choice([
      ignore(string("-")),
      ascii_string([?!..?~], max: 32)
    ])
    |> tag(:msg_id)

  header =
    pri
    |> concat(version)
    |> ignore(string(" "))
    |> concat(datetime)
    |> ignore(string(" "))
    |> concat(hostname)
    |> ignore(string(" "))
    |> concat(app_name)
    |> ignore(string(" "))
    |> concat(proc_id)
    |> ignore(string(" "))
    |> concat(msg_id)

  sd_name = ascii_string([not: ?=, not: ?\s, not: ?], not: ?"], max: 32)

  sd_param =
    sd_name
    |> ignore(string("=\""))
    |> concat(sd_name)
    |> ignore(string("\""))
    |> tag(:sd_param)

  sd_element =
    ignore(string("["))
    |> concat(sd_name |> tag(:sd_id))
    |> repeat(ignore(string(" ")) |> concat(sd_param))
    |> ignore(string("]"))
    |> tag(:sd_element)

  structured_data =
    choice([
      times(sd_element, min: 1),
      string("-")
    ])

  message =
    utf8_string([], min: 1)
    |> eos
    |> tag(:message)

  defparsec(
    :message,
    header
    |> ignore(string(" "))
    |> concat(structured_data)
    |> concat(
      optional(
        ignore(string(" "))
        |> concat(message)
      )
    )
  )
end
