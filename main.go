// contains the main "command" (running) logic
package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/ws")
	return
}
