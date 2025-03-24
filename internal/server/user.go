package server

import (
	"log"
	"net"
	"strings"
)

const (
	CmdHelp   = "help"
	CmdWho    = "who"
	CmdWhoAmI = "whoami"
	CmdRename = "rename"
)

type User struct {
	Name string
	Addr string
	ch   chan string
	conn net.Conn

	server *Server
}

// NewUser creates a new user and starts the message goroutine for this user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		ch:     make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenAndSendMsg()

	return user
}

// Login adds the user to the OnlineMap and broadcasts the Login information to all the users
func (u *User) Login() {
	// add user to online map
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// broadcast user login message
	u.server.Broadcast(u, "Login")

	// welcome the user
	u.ch <- "Welcome to gchat! Please type \"help\" for details."
}

// Logout deletes the user from the OnlineMap and broadcasts the Logout information to all the users
func (u *User) Logout() {
	// remove user from online map
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// broadcast user logout message
	u.server.Broadcast(u, "Logout")
}

// SendMsg sends the message to the current user, the message is sent when SendMsg returns
func (u *User) SendMsg(msg string) {
	msg = strings.TrimSpace(msg) + "\n"
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		log.Println("SendMsg: failed sending message to user, err:" + err.Error())
	}
}

// ProcessMsg processes the message sent from user/client
func (u *User) ProcessMsg(msg string) {
	msg = strings.TrimSpace(msg)

	// skip empty message
	if len(msg) == 0 {
		return
	}

	cmd := strings.Split(msg, " ")[0]
	switch cmd {
	case CmdHelp:
		u.ch <- `Supported commands:
help: show this help message
who: list all online users
whoami: show your username
rename [username]: change your username to [username]
@[user] [message]: send private [message] to the dedicated [user]
[message]: send [message] to everyone online`
		break
	case CmdWho:
		// list all users
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + " is online."
			u.ch <- onlineMsg
		}
		u.server.mapLock.Unlock()
		break
	case CmdWhoAmI:
		u.ch <- u.Name
		break
	case CmdRename:
		l := strings.Split(msg, " ")
		if len(l) != 2 {
			u.ch <- "Username cannot be empty or contain spaces."
			break
		}

		newName := l[1]
		// check if name exists
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.ch <- "Username " + newName + " already exists."
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.Name = newName
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.ch <- "Username has been successfully set to " + u.Name + "."
		}
		break
	default:
		if cmd[:1] == "@" {
			// private chat
			targetUserName := cmd[1:]
			if targetUserName == "" {
				u.ch <- "Target user cannot be empty."
				return
			}
			if targetUserName == u.Name {
				u.ch <- "Hey " + u.Name + "! You are talking to yourself."
				return
			}
			targetUser, ok := u.server.OnlineMap[targetUserName]
			if ok {
				sendMsg := "[" + u.Addr + "]" + u.Name + ": " + msg
				targetUser.ch <- sendMsg
				u.ch <- sendMsg
			} else {
				u.ch <- "User " + targetUserName + " is not online."
			}
		} else {
			// public chat
			u.server.Broadcast(u, msg)
		}
	}
}

// ListenAndSendMsg listens the user channel and send message to the client
func (u *User) ListenAndSendMsg() {
	for {
		msg := <-u.ch
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Println("ListenAndSendMsg: failed sending message to user, err:" + err.Error())
			// when conn.Write returns an error, it usually indicates a network issue and is unrecoverable
			// so we return to quit the goroutine instead of continue
			return
		}
	}
}
