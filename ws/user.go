// Contains logic for front-end WebSocket clients
package ws

import (
	"github.com/gorilla/websocket"
)

// TODO(NilsG-S): could permissions and subscriptions be their own structs?

type User struct {
	email    string
	isAuthed bool
	// Could potentially make this a set...
	permissions   map[string]bool
	subscriptions map[string]bool
}
