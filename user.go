package main

import "net"

type User struct {
	Name string
	Addr string
	ch   chan string
	conn net.Conn
}

// NewUser creates a new user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		ch:   make(chan string),
		conn: conn,
	}

	go user.ListenMessage()

	return user
}

// ListenMessage listens the user channel and send message to the client
func (user *User) ListenMessage() {
	defer user.conn.Close()

	for {
		msg := <-user.ch
		user.conn.Write([]byte(msg + "\n"))
	}
}
