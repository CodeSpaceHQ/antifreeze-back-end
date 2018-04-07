// Contains logic for front-end WebSocket clients
package ws

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/NilsG-S/antifreeze-back-end/common/env"
)

const (
	writeWait = 10 * time.Second
	// This should be just slightly shorter than pingWait
	pingInterval = 30 * time.Second
	// How long the server will wait for a pong response to a ping
	pongWait = 40 * time.Second
)

type user struct {
	server *Server
	key    string
	conn   *websocket.Conn
	send   chan Message
}

func (u *user) writeUser() {
	ticker := time.NewTicker(pingInterval)

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

			err := u.conn.WriteJSON(mes)
			if err != nil {
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
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

// Basically only for authentication
func (u *user) readUser() {
	// unregister if the user disconnects
	defer func() {
		u.server.unregister <- u
		u.conn.Close()
	}()

	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(appData string) error {
		// Only wait so long for pong messages
		u.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		mType, mes, err := u.conn.ReadMessage()
		if err != nil {
			// Handles when client page is closed
			break
		}

		if mType == websocket.TextMessage {
			var uClaims *env.UserClaims
			uClaims, err = u.server.GetAuth().DecodeUser(string(mes))
			if err != nil {
				u.send <- ErrMes{
					Sub:     "/auth",
					Op:      OpError,
					Message: fmt.Sprintf("Auth invalid: %v", err),
				}
				continue
			}

			// Set user key and authorize
			u.key = uClaims.UserKey
			u.server.auth <- u
			u.send <- SuccessMes{
				Sub: "/auth",
				Op:  OpSuccess,
			}
		}
	}

	return
}
