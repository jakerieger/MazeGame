package main

import (
	"fmt"
	"net"
)

const (
	CmdConnect = iota
	CmdDisconnect
	CmdMove
	CmdRequestPosition
)

type Client struct {
	Name  string  `json:"name"`
	Color uint32  `json:"color"`
	PosX  float32 `json:"posX"`
	PosY  float32 `json:"posY"`
}

type ClientConnectRequest struct {
	Name  string `json:"name"`
	Color uint32 `json:"color"`
}

func RegisterClient(name string, color uint32, positionX, positionY float32) *Client {
	return &Client{name, color, positionX, positionY}
}

// { IP : Client{} }
var connectedClients map[string]*Client

func main() {
	// Init our client map
	connectedClients = make(map[string]*Client)

	listener, err := net.Listen("tcp", ":40200")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(listener)

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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	// TODO: Later on, this can be optimized once we know the maximum size of a message block
	buffer := make([]byte, 1024)

	for {
		// The first 4 bytes of the message should contain the command opcode,
		// followed by the json body of the request
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		opcode := int(buffer[0])
		fmt.Println(opcode)

		jsonBody := string(buffer[4:n])
		fmt.Println(jsonBody)

		if opcode == CmdConnect {
			client := RegisterClient("Jake", 0xFFFF0000, 0, 0)
			connectedClients[conn.RemoteAddr().String()] = client
		}

		if opcode == CmdDisconnect {
		}

		if opcode == CmdMove {
		}

		if opcode == CmdRequestPosition {
		}
	}
}
