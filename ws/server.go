package ws

type Server struct {
	// by deviceId to email
	deviceToUser map[int]string
	// by email to all connected user clients
	// when updating, send to all of these
	emailToUsers map[string][]*user
	register     chan *user
	unregister   chan *user
}

// TODO: get users from db on init?

func NewServer() *Server {
	return &Server{
		deviceToUser: make(map[int]*user),
		emailToUsers: make(map[string][]*user),
		register:     make(chan *user),
		unregister:   make(chan *user),
	}
}
