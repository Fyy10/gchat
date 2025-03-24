package main

import (
	"flag"
	"fmt"
	. "gchat/internal/client"
)

var serverIp string
var serverPort int

func init() {
	// usage: ./client -ip 127.0.0.1 -port 8080
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "port", 8080, "server port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	fmt.Printf("Connected to server %s:%d\n", serverIp, serverPort)

	client.Run()
}
