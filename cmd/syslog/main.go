package main

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

type runnerFunc func(string, *lua.LFunction)
type destination interface{}

type listener interface {
	Start() error
	Stop()
}

func main() {
	state := lua.NewState()
	defer state.Close()

	destinations := map[string]destination{}
	listeners := []listener{}

	state.SetGlobal("destination", state.NewFunction(func(state *lua.LState) int {
		name := state.ToString(1)
		configuration := state.ToTable(2)

		destinationType := configuration.RawGetH(lua.LString("type")).(lua.LString)
		destinationConfig := configuration.RawGetH(lua.LString("config")).(*lua.LTable)

		switch destinationType {
		case "file":
			d, err := initializeFileDestination(destinationConfig)
			if err != nil {
				log.Fatalf("error %s", err)
			}
			destinations[name] = d
		default:
			log.Fatalf("unsupported destination type '%s'", destinationType)
		}

		return 0
	}))

	state.SetGlobal("listen", state.NewFunction(func(state *lua.LState) int {
		protocol := state.ToString(1)
		port := state.ToInt(2)
		router := state.ToFunction(3)

		switch protocol {
		case "udp":
			listener, err := initializeUDPlistener(port, router)
			if err != nil {
				log.Fatalf("error %s", err)
			}
			err = listener.Start()
			if err != nil {
				log.Fatalf("error %s", err)
			}
			listeners = append(listeners, listener)
		default:
			log.Fatalf("unsupported listening protocol '%s'", protocol)
		}
		return 0
	}))

	if err := state.DoFile("config.lua"); err != nil {
		panic(err)
	}
}
