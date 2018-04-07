// Contains logic for front-end WebSocket clients
package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type user struct {
	server *Server
	key    string
	conn   *websocket.Conn
	send   chan Message
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

	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, mes, err := u.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		// TODO: Make this applicable stuff other than Auth
		log.Println(mes)
	}

	return
}
