package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

// PING
func ping(args []Value) Value {
	return Value{typ: "string", str: "SEPONG"}
}

// SET
var SETs = map[string]string{}
var SETsMU = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		fmt.Println("Invalid arguments, less than 2")
		return Value{typ: "error", str: "Error argument"}
	}

	key := args[0].bulk
	val := args[1].bulk

	SETsMU.Lock()
	SETs[key] = val
	SETsMU.Unlock()

	return Value{typ: "string", str: "OK"}
}

// GET
func get(args []Value) Value {
	if len(args) != 1 {
		fmt.Println("Invalid arguments, not equal 1")
		return Value{typ: "error", str: "Error argument"}
	}

	key := args[0].bulk

	SETsMU.Lock()
	value, ok := SETs[key]
	SETsMU.Unlock()

	if !ok {
		return Value{typ: "error", str: "Value not found"}
	}

	return Value{typ: "string", str: value}
}
