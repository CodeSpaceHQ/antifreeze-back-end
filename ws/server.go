package ws

type Server struct {
	// by device id
	deviceToUser map[int]*User
	register     chan *User
	unregister   chan *User
}

func newServer() *Server {
	return &Server{
		deviceToUser: make(map[int]*User),
		register:     make(chan *User),
		unregister:   make(chan *User),
	}
}
