package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	port := 6379
	fmt.Printf("Listening on port :%d\n", port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to listen on port %d", port)
		return
	}

	// create database
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer aof.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid value, should be array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid value, array should be more than 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		res := handler(args)

		if command == "SET" {
			aof.Write(value)
		}

		writer.Write(res)
	}
}
