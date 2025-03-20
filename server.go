package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// online user map
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// message broadcasting channel
	Message chan string
}

// NewServer creates a new server based on the given ip and port
func NewServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port, OnlineMap: make(map[string]*User), Message: make(chan string)}
}

// ListenAndBroadcast listens to the server message channel and broadcast the message to all the clients
func (s *Server) ListenAndBroadcast() {
	for {
		msg := <-s.Message

		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.ch <- msg
		}
		s.mapLock.Unlock()
	}
}

// Broadcast sends the user message to the server message channel
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ": " + msg
	s.Message <- sendMsg
}

// Handler reads from connection and handles the requests
func (s *Server) Handler(conn net.Conn) {
	defer conn.Close()
	// TODO
	// create user
	user := NewUser(conn)

	// add user to online map
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// broadcast user online message
	s.Broadcast(user, "Online")

	// read from conn and send message
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.Broadcast(user, "Offline")
				return
			}
			if err != nil && err != io.EOF {
				log.Println("Conn read err:" + err.Error())
				return
			}

			msg := string(buf[:n])
			msg = strings.TrimSpace(msg)
			s.Broadcast(user, msg)
		}
	}()

	// block
	select {}
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

	go s.ListenAndBroadcast()

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
