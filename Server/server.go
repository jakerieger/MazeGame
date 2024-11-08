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
	CmdUpdatePosition
)

type Client struct {
	Name  string  `json:"name"`
	Color uint32  `json:"color"`
	Room  string  `json:"room"`
	PosX  float32 `json:"posX"`
	PosY  float32 `json:"posY"`
	conn  net.Conn
}

type ClientConnectRequest struct {
	Name  string `json:"name"`
	Color uint32 `json:"color"`
	Room  string `json:"roomName"`
}

type ClientPosition struct {
	Name string  `json:"name"`
	PosX float32 `json:"posX"`
	PosY float32 `json:"posY"`
}

func RegisterClient(client *Client) {
	connectedClients[client.Room] = append(connectedClients[client.Room], client)
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

func SendRoomPositions(roomName string) {
	clients := GetRoomClients(roomName)

	// Gather all positions
	var positions []ClientPosition
	for _, client := range clients {
		positions = append(positions, ClientPosition{
			Name: client.Name,
			PosX: client.PosX,
			PosY: client.PosY,
		})
	}

	// Convert positions to JSON
	positionsData, err := json.Marshal(positions)
	if err != nil {
		fmt.Println("Error marshalling positions:", err)
		return
	}

	// Send positions to each client
	for _, client := range clients {
		_, err := client.conn.Write(append([]byte{CmdUpdatePosition}, positionsData...))
		fmt.Printf("Updated client position: (%f, %f)", client.PosX, client.PosY)
		if err != nil {
			fmt.Println("Error sending positions to client:", client.Name, err)
		}
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
	var client *Client

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			if client != nil {
				DisconnectClientFromRoom(client.Room, client.Name)
				SendRoomPositions(client.Room)
			}
			return
		}

		opcode := int(buffer[0])
		jsonBody := string(buffer[4:n])

		switch opcode {
		case CmdConnect:
			data := ClientConnectRequest{}
			err := json.Unmarshal([]byte(jsonBody), &data)
			if err != nil {
				fmt.Println("Error unmarshaling connection request:", err)
				return
			}

			client = &Client{
				Name:  data.Name,
				Color: data.Color,
				Room:  data.Room,
				PosX:  0,
				PosY:  0,
				conn:  conn,
			}
			RegisterClient(client)
			SendRoomPositions(data.Room)

		case CmdDisconnect:
			if client != nil {
				DisconnectClientFromRoom(client.Room, client.Name)
				SendRoomPositions(client.Room)
			}
			return

		case CmdMove:
			if client != nil {
				position := ClientPosition{}
				err := json.Unmarshal([]byte(jsonBody), &position)
				if err != nil {
					fmt.Println("Error unmarshaling position:", err)
					continue
				}

				// Update client position and broadcast it
				client.PosX = position.PosX
				client.PosY = position.PosY
				SendRoomPositions(client.Room)
			}

		case CmdRequestPosition:
			if client != nil {
				SendRoomPositions(client.Room)
			}
		}
	}
}
