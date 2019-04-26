package main

import (
	syslog "github.com/jtarchie/syslog/pkg/log"
	luar "layeh.com/gopher-luar"
	"log"

	lua "github.com/yuin/gopher-lua"
)

type writerFunc func([]byte, *lua.LFunction)
type destination interface{}

type listener interface {
	Start() error
	Stop()
}

type payload struct {
	message []byte
	fun     *lua.LFunction
}

var messages = make(chan payload)

func writer(message []byte, fun *lua.LFunction) {
	messages <- payload{message, fun}
}

func reader(state *lua.LState) {
	log.Println("waiting to read messages from queue")
	for {
		select {
		case p := <-messages:
			log.Println("received message")
			l, _, err := syslog.Parse(p.message)
			if err != nil {
				log.Printf("could not parse message: %s", err)
				continue
			}

			log.Println("calling lua function")
			err = state.CallByParam(lua.P{
				Fn:      p.fun,
				NRet:    1,
				Protect: true,
			}, luar.New(state, l))
			if err != nil {
				log.Printf("could not execute router function: %s", err)
				continue
			}
			ret := state.Get(-1)
			switch ret.(type) {
			case *lua.LNilType:
				log.Printf("drop message on the floor")
			default:
				log.Printf("cannot handle return type: %t", ret)
			}
			state.Pop(1)
		default:
		}
	}
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
			listener, err := initializeUDPlistener(port, router, writer)
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

	reader(state)
}
