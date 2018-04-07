package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Server struct {
	// by user encoding to all connected user clients
	// when updating, send to all of these
	usersByKey map[string]map[*user]bool
	unauthed   map[*user]bool
	register   chan *user
	unregister chan *user
	auth       chan *user

	env.Env
}

// TODO: buffer these channels?
func NewServer(xEnv env.Env) *Server {
	return &Server{
		usersByKey: make(map[string]map[*user]bool),
		unauthed:   make(map[*user]bool),
		register:   make(chan *user),
		unregister: make(chan *user),
		auth:       make(chan *user),
		Env:        xEnv,
	}
}

func (s *Server) RunServer() {
	// Can't use two goroutines because `map` isn't thread safe
	for {
		select {
		case u := <-s.register:
			s.unauthed[u] = true
		case u := <-s.unregister:
			// Authed users
			if _, ok := s.usersByKey[u.key][u]; ok {
				delete(s.usersByKey[u.key], u)
			}

			// Unauthed users
			if _, ok := s.unauthed[u]; ok {
				delete(s.unauthed, u)
			}

			close(u.send)
		case u := <-s.auth:
			// Remove user from unauthed pool
			if _, ok := s.unauthed[u]; ok {
				delete(s.unauthed, u)
			}

			// Initialize user group
			if s.usersByKey[u.key] == nil {
				s.usersByKey[u.key] = make(map[*user]bool)
			}

			// Add user to authed pool
			s.usersByKey[u.key][u] = true
		}
	}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("Unable to upgrade connection: %v", err)
	}

	user := &user{
		server: s,
		conn:   conn,
		send:   make(chan Message, 256),
		// channel of length 256
	}

	s.register <- user

	// Start executing websocket read
	go user.writeUser()
	go user.readUser()

	return nil
}

func (s *Server) PushTemp(userKey, deviceKey string, t env.Temp) {
	mes := TempMes{
		Sub:       "/device/temp",
		Op:        Add,
		DeviceKey: deviceKey,
		Temp:      t.Value,
		Date:      t.Date.Unix(),
	}

	for k, _ := range s.usersByKey[userKey] {
		k.send <- mes
	}
}

// TODO: add functions to inject temperatures, devices, etc. into the server.
