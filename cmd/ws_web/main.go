package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port = 8088
)

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	fmt.Printf("Server running on port %d..\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
