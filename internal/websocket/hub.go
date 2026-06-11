package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Hub struct {
	clients    map[int64]map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	UserID int64       `json:"user_id"`
	Data   interface{} `json:"data"`
}

type Client struct {
	Conn   *websocket.Conn
	UserID int64
	Hub    *Hub
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; !ok {
				h.clients[client.UserID] = make(map[*websocket.Conn]bool)
			}
			h.clients[client.UserID][client.Conn] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.UserID]; ok {
				delete(conns, client.Conn)
				if len(conns) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.RLock()
			conns := h.clients[msg.UserID]
			data, _ := json.Marshal(msg.Data)
			for conn := range conns {
				err := conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					conn.Close()
					delete(conns, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID int64, data interface{}) {
	h.broadcast <- Message{UserID: userID, Data: data}
}
