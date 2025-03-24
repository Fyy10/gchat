package main

import (
	"flag"
	. "gchat/internal/server"
)

var serverIp string
var serverPort int

func init() {
	// usage: ./server -ip 0.0.0.0 -port 8080
	flag.StringVar(&serverIp, "ip", "0.0.0.0", "server ip")
	flag.IntVar(&serverPort, "port", 8080, "server port")
}

func main() {
	flag.Parse()
	server := NewServer(serverIp, serverPort)
	server.Run()
}
