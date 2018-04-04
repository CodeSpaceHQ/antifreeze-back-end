package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	// "github.com/NilsG-S/antifreeze-back-end/common/device"
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
	// usersByKey   map[string]map[*user]bool
	register   chan *user
	unregister chan *user
}

func NewServer() *Server {
	// TODO: get users from db on init?

	return &Server{
		deviceToUser: make(map[int]string),
		emailToUsers: make(map[string]map[*user]bool),
		register:     make(chan *user),
		unregister:   make(chan *user),
	}
}

func (s *Server) RunServer() {
	// Can't use two goroutines because `map` isn't thread safe
	for {
		select {
		case u := <-s.register:
			if s.emailToUsers[u.email] == nil {
				s.emailToUsers[u.email] = make(map[*user]bool)
			}

			s.emailToUsers[u.email][u] = true
		case user := <-s.unregister:
			if _, ok := s.emailToUsers[user.email][user]; ok {
				delete(s.emailToUsers[user.email], user)
				close(user.send)
			}
		}
	}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	user := &user{
		server: s,
		// TODO: This is a stopgap. Replace with authentication
		email: "test@ttu.edu",
		perms: unauthed,
		conn:  conn,
		// channel of length 256
		send: make(chan Message, 256),
	}

	s.register <- user

	// Start executing websocket read
	go user.writeUser()
	go user.readUser()
}

func (s *Server) POSTDeviceHistory(mes Temperature) {
	id := mes.DeviceID
	email, ok := s.deviceToUser[id]
	if !ok {
		return
	}

	for k, _ := range s.emailToUsers[email] {
		k.send <- mes
	}
}

// func (s *Server) PushTemp(userKey, deviceKey string, t device.Temp) {
// 	// TODO: move message to WS
// 	return
// }

func (s *Server) POSTUserDevices(id int, email string) {
	s.deviceToUser[id] = email
}

// Add functions to inject temperatures, devices, etc. into the server.
