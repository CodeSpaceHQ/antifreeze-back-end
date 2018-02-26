// contains the main "command" (running) logic
package main

import (
	// "net/http"
	"fmt"
)

func main() {
	// http.HandleFunc("/ws")

	var test []string
	fmt.Println(test == nil)
	test = append(test, "stuff")
	fmt.Println(test == nil)

	var test2 map[int]string
	fmt.Println(test2 == nil)
	test2 = make(map[int]string)
	fmt.Println(test2 == nil)

	return
}
