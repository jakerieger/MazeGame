package main

import (
	"fmt"
	"net"
)

// { IP : Client{} }
var connectedClients map[string]*Client

func main() {
	listener, err := net.Listen("tcp", ":40200")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 40200")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// TODO: Later on, this can be optimized once we know the maximum size of a message block
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received: %s", buffer[:n])
	}
}
