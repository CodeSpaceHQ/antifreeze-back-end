package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
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
	// by user encoding to all connected user clients
	// when updating, send to all of these
	usersByKey map[string]map[*user]bool
	register   chan *user
	unregister chan *user
}

// TODO: This function may not be necessary
func NewServer() *Server {
	return &Server{
		usersByKey: make(map[string]map[*user]bool),
		register:   make(chan *user),
		unregister: make(chan *user),
	}
}

func (s *Server) RunServer() {
	// Can't use two goroutines because `map` isn't thread safe
	for {
		select {
		case u := <-s.register:
			if s.usersByKey[u.key] == nil {
				s.usersByKey[u.key] = make(map[*user]bool)
			}

			s.usersByKey[u.key][u] = true
		case user := <-s.unregister:
			if _, ok := s.usersByKey[user.key][user]; ok {
				delete(s.usersByKey[user.key], user)
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
		// For now, you can only get this value by printing the contents of the device JWT
		key:   "Eg8KBFVzZXIQgICAgIDkkQo",
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

func (s *Server) PushTemp(userKey, deviceKey string, t env.Temp) {
	mes := Temperature{
		Sub:       "/device/history",
		Op:        Add,
		DeviceKey: deviceKey,
		Temp:      t.Value,
		Date:      t.Date.Unix(),
	}

	for k, _ := range s.usersByKey[userKey] {
		k.send <- mes
	}
}

// Add functions to inject temperatures, devices, etc. into the server.
