package main

import (
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// hub 의 run 은 hub 의 여러 채널로 넘어오는 데이터를 처리한다.
// register 채널로 넘어오는 client 를 등록하고,
// unregister 채널로 넘어오는 client 를 삭제하고
// broadcast 채널로 넘어오는 message 를 전체 client 에 뿌린다.
// broadcast 채넘로 메시지가 수신이 확인되는 경우, client 를 순회하며, client 의 send 채널로 message 를 전달한다.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Register client: %s", client.conn.RemoteAddr())
			// buffer 가 없는 채널로 설정된 broadcast 는 message 가 개시 된 후 채널에서 수신하지 않으면 블록된다.
			// broadcast 채널의 수신 부분이 register 하는 부분과 같은 함수인 run() 함수에서 select 문으로 처리가 되고 있기 때문에,
			// broadcast 채널로 메시지를 여기서 전송하는 것은 고루틴을 블록시킨다.
			// h.broadcast <- []byte(fmt.Sprintf("%s has joined the chat room", client.conn.RemoteAddr().String()))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Printf("Broadcast: %s", message)
			for client := range h.clients {
				select {
				case client.send <- message:
					log.Printf("Sending to client: %s", client.conn.RemoteAddr())
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}
