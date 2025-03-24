package client

import (
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

// ListenAndPrintMsg reads message from conn and prints the message to user console
func (c *Client) ListenAndPrintMsg() {
	_, err := io.Copy(os.Stdout, c.conn)
	fmt.Println("Connection closed by server")
	if err != nil {
		fmt.Println("Failed receiving message from server, err:", err)
	}
}

func (c *Client) Run() {
	go c.ListenAndPrintMsg()

	// FIXME: when connection is closed by server, io.Copy does not return, until user sends some message with TWO returns
	_, err := io.Copy(c.conn, os.Stdin)
	if err != nil {
		fmt.Println("Failed sending message to server, err:", err)
	}
}
