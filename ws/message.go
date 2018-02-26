package ws

const (
	add    int = 1
	remove int = 2
	update int = 3
)

// support multiple message types
type message interface {
	getSub() string
	getOp() int
}
