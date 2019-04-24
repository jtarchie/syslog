package main

import (
	"fmt"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

type fileDestination struct {
	path            string
	messageModifier *lua.LFunction
}

func initializeFileDestination(config *lua.LTable) (*fileDestination, error) {
	d := &fileDestination{}

	path := config.RawGetH(lua.LString("path")).(lua.LString).String()
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("cannot expand path '%s' to absolute path", path)
	}

	messageModifierFn, ok := config.RawGetH(lua.LString("message")).(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("message modifier method is not a valid lua function")
	}

	d.path = fullPath
	d.messageModifier = messageModifierFn

	return d, nil
}
