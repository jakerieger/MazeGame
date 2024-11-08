package main

type Client struct {
	name      string
	color     uint32 // RGBA
	positionX float32
	positionY float32
}

func RegisterClient(name string, color uint32, positionX, positionY float32) *Client {
	return &Client{name, color, positionX, positionY}
}
