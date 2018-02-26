package ws

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

			// Start executing websocket read
			go user.writeUser()
			go user.readUser()
		case user := <-v.unregister:
			if _, ok := v.emailToUsers[user.email][user]; ok {
				delete(v.emailToUsers[user.email], user)
				close(user.send)
			}
		}
	}
}

// Functions to register new users

func (s *Server) POSTDeviceHistory(mes temp) {
	email := s.deviceToUser[mes.deviceId]

	for k, _ := range s.emailToUsers[email] {
		k.send <- mes
	}
}

// Add functions to inject temperatures, devices, etc. into the server.
