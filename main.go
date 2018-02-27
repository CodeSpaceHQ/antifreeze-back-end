// contains the main "command" (running) logic
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/NilsG-S/antifreeze-back-end/ws"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	server := ws.NewServer()
	go server.RunServer()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.Register(w, r)
		}
	})

	http.HandleFunc("/device/history", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.POSTDeviceHistory(w, r)
		}
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}
