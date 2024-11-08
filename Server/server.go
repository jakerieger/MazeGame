package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
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
	Name     string `json:"name"`
	Color    uint32 `json:"color"`
	RoomName string `json:"roomName"`
}

type ClientPosition struct {
	PosX float32 `json:"posX"`
	PosY float32 `json:"posY"`
}

func RegisterClient(roomName string, name string, color uint32, positionX, positionY float32) {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()

	client := &Client{name, color, positionX, positionY}
	connectedClients[roomName] = append(connectedClients[roomName], client)
}

func GetRoomClients(roomName string) []*Client {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()

	return connectedClients[roomName]
}

func DisconnectClientFromRoom(roomName string, clientName string) {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()

	var newClients []*Client
	for _, client := range connectedClients[roomName] {
		if client.Name != clientName {
			newClients = append(newClients, client)
		}
	}

	// Delete room if no players are present
	if len(newClients) == 0 {
		delete(connectedClients, roomName)
	} else {
		connectedClients[roomName] = newClients
	}
}

var (
	connectedClients = make(map[string][]*Client)
	clientMapMutex   = sync.RWMutex{}
)

func main() {
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
			data := ClientConnectRequest{}
			err := json.Unmarshal([]byte(jsonBody), &data)
			if err != nil {
				fmt.Println(err)
				return
			}
			RegisterClient(data.RoomName, data.Name, data.Color, 0, 0)
		}

		if opcode == CmdDisconnect {
			fmt.Println("Client disconnected.")
		}

		if opcode == CmdMove {
		}

		if opcode == CmdRequestPosition {
		}
	}
}
