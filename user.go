package main

import (
	"log"
	"net"
)

type User struct {
	Name string
	Addr string
	ch   chan string
	conn net.Conn

	server *Server
}

// NewUser creates a new user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		ch:     make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (u *User) Login() {
	// add user to online map
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// broadcast user login message
	u.server.Broadcast(u, "Login")
}

func (u *User) Logout() {
	// remove user from online map
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// broadcast user logout message
	u.server.Broadcast(u, "Logout")
}

func (u *User) Send(msg string) {
	switch msg {
	case "who":
		// list all users
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": I am online."
			u.ch <- onlineMsg
		}
	case "whoami":
		u.ch <- u.Name
		u.server.mapLock.Unlock()
		break
	default:
		u.server.Broadcast(u, msg)
	}
}

// ListenMessage listens the user channel and send message to the client
func (u *User) ListenMessage() {
	for {
		msg := <-u.ch
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Println("Failed writing to socket, err:" + err.Error())
			continue
		}
	}
}
