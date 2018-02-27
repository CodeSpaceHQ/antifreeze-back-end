package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	// TODO: This might not be relevant for JSON
	maxMessageSize = 512
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Server struct {
	// by deviceId to email
	deviceToUser map[int]string
	// by email to all connected user clients
	// when updating, send to all of these
	emailToUsers map[string]map[*user]bool
	register     chan *user
	unregister   chan *user
}

func NewServer() *Server {
	// TODO: get users from db on init?

	return &Server{
		deviceToUser: make(map[int]*user),
		emailToUsers: make(map[string]map[*user]bool),
		register:     make(chan *user),
		unregister:   make(chan *user),
	}
}

func (v *Server) RunServer() {
	// Can't use two goroutines because `map` isn't thread safe
	for {
		select {
		case user := <-v.register:
			if v.emailToUsers[user.email] == nil {
				v.emailToUsers[user.email] = make(map[*user]bool)
			}

			v.emailToUsers[user.email][user] = true
		case user := <-v.unregister:
			if _, ok := v.emailToUsers[user.email][user]; ok {
				delete(v.emailToUsers[user.email], user)
				close(user.send)
			}
		}
	}
}

func (s *Server) Register(w http.ResponseWriter, r http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: ensure this value retrieval works
	user := &user{
		email: r.FormValue("email"),
		perms: unauthed,
		conn:  conn,
		// channel of length 256
		send: make(chan []mes, 256),
	}

	s.register <- user

	// Start executing websocket read
	go user.writeUser()
	go user.readUser()
}

func (s *Server) POSTDeviceHistory(w http.ResponseWriter, r http.Request) {
	id := r.FormValue("deviceId")
	email, ok := s.deviceToUser[id]
	if !ok {
		return
	}

	for k, _ := range s.emailToUsers[email] {
		k.send <- &temp{
			sub:      "/device/history",
			op:       add,
			deviceId: id,
			temp:     r.FormValue("temp"),
			time:     r.FormValue("time"),
		}
	}
}

// Add functions to inject temperatures, devices, etc. into the server.
