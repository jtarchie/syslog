validations = %{
  "[no] empty input" => fn -> {:error, _} = Syslog.parse("") end,
  # "[no] multiple syslog messages on multiple lines" => fn ->
  #   {:error, _} = Syslog.parse("<1>1 - - - - - -\n<2>1 - - - - - -")
  # end,
  "[no] malformed structured data" => fn -> {:error, _} = Syslog.parse(~S(<1>1 - - - - - X)) end,
  "[ok] minimal" => fn -> {:ok, _} = Syslog.parse(~S(<1>1 - - - - - -)) end,
  "[ok] average message" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<29>1 2016-02-21T04:32:57+00:00 web1 someservice - - [origin x-service="someservice"][meta sequenceId="14125553"] 127.0.0.1 - - 1456029177 "GET /v1/ok HTTP/1.1" 200 145 "-" "hacheck 0.9.0" 24306 127.0.0.1:40124 575)
      )
  end,
  "[ok] complicated message" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<78>1 2016-01-15T00:04:01Z host1 CROND 10391 - [meta sequenceId="29" sequenceBlah="foo"][my key="value"] some_message)
      )
  end,
  "[ok] very long message" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<190>1 2016-02-21T01:19:11+00:00 batch6sj - - - [meta sequenceId="21881798" x-group="37051387"][origin x-service="tracking"] metascutellar conversationalist nephralgic exogenetic graphy streng outtaken acouasm amateurism prenotice Lyonese bedull antigrammatical diosphenol gastriloquial bayoneteer sweetener naggy roughhouser dighter addend sulphacid uneffectless ferroprussiate reveal Mazdaist plaudite Australasian distributival wiseman rumness Seidel topazine shahdom sinsion mesmerically pinguedinous ophthalmotonometer scuppler wound eciliate expectedly carriwitchet dictatorialism bindweb pyelitic idic atule kokoon poultryproof rusticial seedlip nitrosate splenadenoma holobenthic uneternal Phocaean epigenic doubtlessly indirection torticollar robomb adoptedly outspeak wappenschawing talalgia Goop domitic savola unstrafed carded unmagnified mythologically orchester obliteration imperialine undisobeyed galvanoplastical cycloplegia quinquennia foremean umbonal marcgraviaceous happenstance theoretical necropoles wayworn Igbira pseudoangelic raising unfrounced lamasary centaurial Japanolatry microlepidoptera)
      )
  end,
  "[ok] all max length and complete" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<191>999 2018-12-31T23:59:59.999999-23:59 abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabc abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdef abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzab abcdefghilmnopqrstuvzabcdefghilm [an@id key1="val1" key2="val2"][another@id key1="val1"] Some message "GET")
      )
  end,
  "[ok] all max length except structured data and message" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<191>999 2018-12-31T23:59:59.999999-23:59 abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabc abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdef abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzab abcdefghilmnopqrstuvzabcdefghilm -)
      )
  end,
  "[ok] minimal with message containing newline" => fn ->
    {:ok, _} = Syslog.parse(~S(<1>1 - - - - - - x\x0Ay))
  end,
  "[ok] w/o procid, w/o structured data, with message starting with BOM" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - \xEF\xBB\xBF'su root' failed for lonvick on /dev/pts/8)
      )
  end,
  "[ok] minimal with UTF-8 message" => fn ->
    {:ok, _} = Syslog.parse(~S(<0>1 - - - - - - ⠊⠀⠉⠁⠝⠀⠑⠁⠞⠀⠛⠇⠁⠎⠎⠀⠁⠝⠙⠀⠊⠞⠀⠙⠕⠑⠎⠝⠞⠀⠓⠥⠗⠞⠀⠍⠑))
  end,
  "[ok] with structured data id, w/o structured data params" => fn ->
    {:ok, _} = Syslog.parse(~S(<29>50 2016-01-15T01:00:43Z hn S - - [my@id]))
  end,
  "[ok] with multiple structured data" => fn ->
    {:ok, _} =
      Syslog.parse(~S(<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 k="v"][my@id2 c="val"]))
  end,
  "[ok] with escaped backslash within structured data param value, with message" => fn ->
    {:ok, _} =
      Syslog.parse(~S(<29>50 2016-01-15T01:00:43Z hn S - - [meta es="\\valid"] 1452819643))
  end,
  "[ok] with UTF-8 structured data param value, with message" => fn ->
    {:ok, _} =
      Syslog.parse(
        ~S(<78>1 2016-01-15T00:04:01+00:00 host1 CROND 10391 - [sdid x="⌘"] some_message)
      )
  end
}

validations
|> Enum.each(fn {name, func} ->
  IO.puts("testing validation of '#{name}'")
  func.()
end)

require ExProf.Macro
ExProf.Macro.profile do
  {:ok, _} =
    Syslog.parse(
      ~S(<29>1 2016-02-21T04:32:57+00:00 web1 someservice - - [origin x-service="someservice"][meta sequenceId="14125553"] 127.0.0.1 - - 1456029177 "GET /v1/ok HTTP/1.1" 200 145 "-" "hacheck 0.9.0" 24306 127.0.0.1:40124 575)
    )
end

Benchee.run(validations,
  print: [fast_warning: false],
  console: [
    extended_statistics: [:sample_size]
  ]
)
