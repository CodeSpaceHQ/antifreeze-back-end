// contains the main "command" (running) logic
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/NilsG-S/antifreeze-back-end/common"
	"github.com/NilsG-S/antifreeze-back-end/ws"
)

var addr = flag.String("addr", ":8080", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	server := ws.NewServer()
	go server.RunServer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			http.ServeFile(w, r, "home.html")
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// server.Register(w, r)
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(conn)
		}
	})

	http.HandleFunc("/device/history", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			mes := common.Temperature{
				Sub:      "/device/history",
				Op:       common.Add,
				DeviceID: 1,
				Temp:     32,
				Time:     time.Now(),
			}

			server.POSTDeviceHistory(mes)
		}
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}
