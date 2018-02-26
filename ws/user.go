// Contains logic for front-end WebSocket clients
package ws

import (
	"github.com/gorilla/websocket"
)

// TODO(NilsG-S): could permissions and subscriptions be their own structs?

type user struct {
	email string
	// Could potentially make this a set...
	permissions map[string]bool
	// used to decide whether to send information
	// technically not necessay under the current proposal
	subscriptions map[string]bool
}
