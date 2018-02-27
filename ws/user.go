// Contains logic for front-end WebSocket clients
package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

// TODO(NilsG-S): could permissions and subscriptions be their own structs?

type perms struct {
	authed bool
}

var (
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
	perms *perms
	// used to decide whether to send information
	// technically not necessay under the current proposal
	// subs map[string]bool
	conn *websocket.Conn
	send chan message
}

func (u *user) writeUser() {
	// TODO: convert this to use JSON
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		u.conn.Close()
	}()

	for {
		select {
		case mes, ok := <-u.send:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// If someone has unregistered this user (which closes the channel)
			if !ok {
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(mes.getSub()))

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// TODO: will these ever need to send to other users?
func (u *user) readUser() {
	// unregister if the user disconnects
	return
}

// ensures a given user has the right permissions
func (v *user) checkAuth(req *perms) bool {
	// TODO: do stuff
	return true
}
