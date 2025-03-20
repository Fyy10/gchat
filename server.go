package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer creates a new server based on the given ip and port
func NewServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port}
}

// Handler reads from connection and handles the requests
func (s *Server) Handler(conn net.Conn) {
	defer conn.Close()
	// TODO
	conn.Write([]byte("Hello World!"))
}

// Run starts the server and listen to the socket
func (s *Server) Run() {
	// listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		panic("net.Listen failed, err: " + err.Error())
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept failed, err: " + err.Error())
			continue
		}

		// do handler
		go s.Handler(conn)
	}
}
