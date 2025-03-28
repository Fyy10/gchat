package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	Timeout = 600 // seconds
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
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
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

// Broadcast sends the user/client message to the server message channel
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	s.Message <- sendMsg
}

// Handler reads from connection and handles the requests
func (s *Server) Handler(conn net.Conn) {
	// create user
	user := NewUser(conn, s)

	user.Login()

	isActive := make(chan bool)

	// read from conn and process messages
	go func() {
		defer user.Logout()
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				return
			}
			if err != nil && err != io.EOF {
				log.Println("Conn read err:" + err.Error())
				return
			}

			msg := string(buf[:n])
			user.ProcessMsg(msg)
			// any msg from user indicates that the user is active
			isActive <- true
		}
	}()

	// set timer
	for {
		select {
		case <-isActive:
			break
		case <-time.After(time.Second * Timeout):
			// timeout
			// kick out the user
			user.SendMsg("You are kicked out for being inactive.")
			err := conn.Close()
			if err != nil {
				log.Printf("Error closing connection: %v", err)
			}
			// after the user is kicked out, the current handler function should also return
			return
		}
	}
}

// Run starts the server and listen to the socket
func (s *Server) Run() {
	// listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Fatalln("net.Listen failed, err: " + err.Error())
	}
	// close listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}(listener)

	go s.ListenAndBroadcast()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept failed, err: " + err.Error())
			continue
		}
		log.Println("Established connection from " + conn.RemoteAddr().String())

		// do handler
		go s.Handler(conn)
	}
}
