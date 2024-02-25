package main

import (
	"flag"
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

func handleConnection(conn net.Conn, storage Storage, repl Replication) {
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
		case "set":
			key := commands[1]
			value := commands[2]
			var expireIn *int64
			if len(commands) >= 5 {
				command := commands[3]
				switch command {
				case "px":
					howLongMs := commands[4]
					if val, err := strconv.ParseInt(howLongMs, 10, 64); err == nil {
						expireIn = &val
					} else {
						fmt.Errorf("Error px ParseInt :", err.Error())
						os.Exit(1)
					}
				}
			}
			storage.Set(key, value, expireIn)
			ok(conn)
		case "get":
			key := commands[1]
			value := get(storage, key)
			echo(conn, value)
		case "info":
			whatInfo := commands[1]
			switch whatInfo {
			case "replication":
				echo(conn, fmt.Sprintf("role:%s", repl.GetStatus().Role))
			}
		}
	}
}

func ok(conn net.Conn) {
	msg := fmt.Sprintf("+OK\r\n")
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Errorf("Error responding ok:", err.Error())
		os.Exit(1)
	}
}

func get(storage Storage, key string) string {
	return storage.Get(key)
}

func echo(conn net.Conn, message string) {
	msg := fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
	if message == "-1" { // Expired
		msg = fmt.Sprintf("$%s\r\n", message)
	}
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
	port := flag.String("port", "6379", "port number")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", *port))
	if err != nil {
		fmt.Errorf("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	redis := NewMemoryStorage()
	repl := NewReplicationStorage()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, redis, repl)
	}
}
