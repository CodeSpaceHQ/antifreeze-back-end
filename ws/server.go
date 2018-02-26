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

// TODO: get users from db on init?

func NewServer() *Server {
	return &Server{
		deviceToUser: make(map[int]*user),
		emailToUsers: make(map[string]map[*user]bool),
		register:     make(chan *user),
		unregister:   make(chan *user),
	}
}

func (v *Server) Run() {
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

func (v *Server) register(user *user) error {
	if v.emailToUsers[user.email] == nil {
		v.emailToUsers[user.email] = make(map[*user]bool)
	}

	v.emailToUsers[user.email][user] = true
}

func (v *Server) unregister(user *user) error {
	if _, ok := v.emailToUsers[user.email][user]; ok {
		delete(v.emailToUsers[user.email], user)
		close(user.send)
	}
}

// Add functions to inject temperatures, devices, etc. into the server.
