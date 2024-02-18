package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./app <command>")
	}

	if args[1] == "ping" {
		fmt.Print("pong")
		os.Exit(0)
	} else {
		fmt.Println("Unknown command:", args[1])
		os.Exit(1)
	}

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
