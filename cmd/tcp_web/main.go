package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	client    = make(map[net.Conn]bool)
	broadcast = make(chan []byte)
	mutex     = sync.Mutex{}
)

const (
	port = 8081
)

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	defer listen.Close()

	fmt.Printf("Server is running on %d...\n", port)
	go handleBroadcast()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		mutex.Lock()
		client[conn] = true
		mutex.Unlock()
		go handleConnection(conn)
		broadcast <- []byte(fmt.Sprintf("%s has joined", conn.RemoteAddr().String()))
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		broadcast <- []byte(fmt.Sprintf("%s: %s", conn.RemoteAddr().String(), msg))
	}

	mutex.Lock()
	delete(client, conn)
	mutex.Unlock()
	broadcast <- []byte(fmt.Sprintf("%s has left", conn.RemoteAddr().String()))
}

func handleBroadcast() {
	for {
		select {
		case msg := <-broadcast:
			mutex.Lock()
			for conn := range client {
				fmt.Fprintln(conn, string(msg))
			}
			mutex.Unlock()
		}
	}
}
