package main

import (
	lua "github.com/yuin/gopher-lua"
	"log"
	"path/filepath"
)

type fileDestination struct {
	path string
	messageModifier *lua.LFunction
}

type destination interface {}

func main() {
	state := lua.NewState()
	defer state.Close()

	destinations := map[string]destination{}

	state.SetGlobal("destination", state.NewFunction(func(state *lua.LState) int {
		name := state.ToString(1)
		configuration := state.ToTable(2)

		destinationType := configuration.RawGetH(lua.LString("type")).(lua.LString)
		destinationConfig := configuration.RawGetH(lua.LString("config")).(*lua.LTable)

		switch destinationType {
		case "file":
			path := destinationConfig.RawGetH(lua.LString("path")).(lua.LString).String()
			fullPath, err := filepath.Abs(path)
			if err != nil {
				log.Fatalf("cannot expand path '%s' to absolute path", path)
			}

			messageModifierFn, ok := destinationConfig.RawGetH(lua.LString("path")).(*lua.LFunction)
			if !ok {
				log.Fatalf("message modifier method is not a valid lua function")
			}
			destinations[name] = fileDestination{
				path: fullPath,
				messageModifier: messageModifierFn,
			}
		default:
			log.Fatalf("unsupported destination type '%s'", destinationType)
		}

		return 0
	}))

	state.SetGlobal("listen", state.NewFunction(func(lState *lua.LState) int {
		return 0
	}))

	if err := state.DoFile("config.lua"); err != nil {
		panic(err)
	}
}