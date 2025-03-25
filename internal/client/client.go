package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	conn       net.Conn
}

// NewClient creates a new client and connects to the server
func NewClient(serverIp string, serverPort int) *Client {
	// create client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		log.Fatalln("net.Dial failed, err:", err)
	}

	client.conn = conn

	// return
	return client
}

// ListenAndSend reads message from user and sends the message to the server
func (c *Client) ListenAndSend() {
	defer c.conn.Close()

	_, err := io.Copy(c.conn, os.Stdin)
	if err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Println("Failed sending message to server, err:", err)
	}
}

func (c *Client) Run() {
	defer c.conn.Close()

	go c.ListenAndSend()

	// read message from conn and print the message to user console
	_, err := io.Copy(os.Stdout, c.conn)
	if err != nil && !errors.Is(err, net.ErrClosed) {
		fmt.Println("Failed receiving message from server, err:", err)
	}
}
