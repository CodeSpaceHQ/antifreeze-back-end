// Contains logic for front-end WebSocket clients
package ws

import (
	"github.com/gorilla/websocket"
)

// TODO(NilsG-S): could permissions and subscriptions be their own structs?

type perms struct {
	authed bool
}

const (
	authed *perms = &perms{
		authed: true,
	}
	unauthed *perms = &perms{
		authed: false,
	}
)

// type subs struct {
// 	userDevices   bool
// 	devicesAlarm  bool
// 	deviceHistory bool
// }

type user struct {
	email string
	perms map[string]bool
	// used to decide whether to send information
	// technically not necessay under the current proposal
	// subs map[string]bool
	conn *websocket.Conn
	send chan []message
}

// ensures a given user has the right permissions
func (v *user) checkAuth(req *perms) bool {
	// TODO: do stuff
	return true
}
