package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func parseRedisMessage(message string) ([]string, error) {
	var parts []string

	segments := strings.Split(message, "\r\n")

	for i := int(1); i < len(segments)-1; i += 2 {
		length, err := strconv.Atoi(segments[i][1:])
		if err != nil {
			return nil, err
		}

		parts = append(parts, segments[i+1][:length])
	}
	return parts, nil
}

func handleConnection(conn net.Conn) {
	fmt.Println("Handling new connection")

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Errorf("Error reading from connection:", err.Error())
				os.Exit(1)
			}
			return
		}
		fmt.Printf("received %d bytes\n", n)
		fmt.Printf("received the following data: %s", string(buf[:n]))

		commands, err := parseRedisMessage(string(buf[:n]))
		if err != nil {
			fmt.Errorf("Error parseRedisMessage:", err.Error())
			os.Exit(1)
		}
		fmt.Println("commands", commands)

		switch commands[0] {
		case "ping":
			pong(conn)
		case "echo":
			echo(conn, commands[1])
		}
	}
}

func echo(conn net.Conn, message string) {
	msg := fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
	n, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Errorf("Error responding to echo:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("sent %d bytes\n", n)
	fmt.Printf("sent the following data: %s", msg)
}

func pong(conn net.Conn) {
	pong := []byte("+PONG\r\n")
	n, err := conn.Write(pong)
	if err != nil {
		fmt.Errorf("Error responding to ping:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("sent %d bytes\n", n)
	fmt.Printf("sent the following data: %s", string(pong))
}

func isPing(command string) bool {
	return command == "ping"
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Errorf("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
