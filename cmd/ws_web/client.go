package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.hub.broadcast <- []byte(fmt.Sprintf("%s has left.", c.conn.RemoteAddr().String()))
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// 서버에서 전송되는 ping 에 의해 브라우저에 내장된 WebSocket API 가 자동으로 pong 응답을 보낸다.
	// pong 핸들러는 클라이언트의 pong 을 받아 연결이 정상임을 확인하고, read 함수에서 읽기 제한 시간을 갱신해 타임아웃을 관리한다.
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// NextWriter 가 호출될떄 마다 새로운 프레임이 생기기 때문에 비용이 발생
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 한번의 NextWriter 를 생성해서 여러 메시지를 병합해 처리하여 자원 오버헤드를 줄인다.
			// 네트워크 호출을 줄여 성능을 개선할 수 있다.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// ping message 를 주기적으로 client 에 보내 client 의 상태를 확인한다.
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) join() {
	c.hub.register <- c
	c.hub.broadcast <- []byte(fmt.Sprintf("%s has joined.", c.conn.RemoteAddr().String()))
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.join()

	go client.write()
	go client.read()
}
