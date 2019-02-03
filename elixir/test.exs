defmodule Bug do
    import NimbleParsec
    defparsec :parse, integer(min: 1, max: 3)
end

'"' == Bug.parse("34")
[1,2,3] == Bug.parse("123")