package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Hub struct {
	clients    map[int64]map[*websocket.Conn]bool // userID -> connections
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	Conn   *websocket.Conn
	UserID int64
	Hub    *Hub
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
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
			var payload struct {
				UserID int64 `json:"user_id"`
				Data   json.RawMessage
			}
			if err := json.Unmarshal(msg, &payload); err != nil {
				continue
			}
			h.mu.RLock()
			conns := h.clients[payload.UserID]
			for conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, payload.Data); err != nil {
					conn.Close()
					delete(conns, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID int64, data interface{}) {
	msg, _ := json.Marshal(data)
	h.broadcast <- msg
}
