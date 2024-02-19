package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	fmt.Println("Handling new connection")

	for {

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err.Error())
				os.Exit(1)
			}
			return
		}
		fmt.Printf("received %d bytes\n", n)
		fmt.Printf("received the following data: %s", string(buf[:n]))
		if isPing(buf[:n]) {
			pong(conn)
		}
	}
}

func pong(conn net.Conn) {
	pong := []byte("+PONG\r\n")
	n, err := conn.Write(pong)
	if err != nil {
		fmt.Println("Error responsding to ping:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("sent %d bytes\n", n)
	fmt.Printf("sent the following data: %s", string(pong))
}

func isPing(command []byte) bool {
	ping := "*1\r\n$4\r\nping\r\n" // RESP
	if string(command) == ping {
		return true
	}
	return false
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
