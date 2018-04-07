// Contains logic for front-end WebSocket clients
package ws

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
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
	// TODO: convert this to use JSON
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

			w, err := u.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(mes.GetSub()))

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
		_, mes, err := u.conn.ReadMessage()
		if err != nil {
			// Handles when client page is closed
			break
		}

		fmt.Println(mes)
	}

	return
}
