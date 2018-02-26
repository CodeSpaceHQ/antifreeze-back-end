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

func (u *user) writeUser() {
	// TODO: convert this to use JSON
	for {
		mes, ok := <-u.send

		// If someone has unregistered this user (which closes the channel)
		if !ok {
			return
		}

		w, err := u.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(mes)
	}
}

func (u *user) readUser() {
	// unregister if the user disconnects
	return
}

// ensures a given user has the right permissions
func (v *user) checkAuth(req *perms) bool {
	// TODO: do stuff
	return true
}
