package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}
)

type EventType string

const (
	EventTypeNewRequest  EventType = "NEW_REQUEST"
	EventTypeNewOffer    EventType = "NEW_OFFER"
	EventTypeOfferAccept EventType = "OFFER_ACCEPTED"
)

type WSEvent struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client holds a WebSocket connection and a set of category IDs the expert is subscribed to
type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	categoryIDs []uuid.UUID // categories this expert cares about (empty = receive all)
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastEvent sends a WS event to ALL connected clients (used for OFFER_ACCEPTED etc.)
func (h *Hub) BroadcastEvent(eventType EventType, payload interface{}) {
	event := WSEvent{Type: eventType, Payload: payload}
	msg, err := json.Marshal(event)
	if err != nil {
		log.Println("Error marshalling WS event:", err)
		return
	}
	h.broadcast <- msg
}

// BroadcastToCategory sends a WS event only to clients subscribed to the given categoryID.
// Clients with no categories registered receive it too (backward compat / non-expert clients).
func (h *Hub) BroadcastToCategory(categoryID uuid.UUID, eventType EventType, payload interface{}) {
	event := WSEvent{Type: eventType, Payload: payload}
	msg, err := json.Marshal(event)
	if err != nil {
		log.Println("Error marshalling WS event:", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		if matchesCategory(client.categoryIDs, categoryID) {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// matchesCategory returns true if the client has no filter (empty list) or if categoryID is in the list.
func matchesCategory(clientCats []uuid.UUID, categoryID uuid.UUID) bool {
	if len(clientCats) == 0 {
		return true // no filter set — receive all (e.g., user clients)
	}
	for _, c := range clientCats {
		if c == categoryID {
			return true
		}
	}
	return false
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleConnection upgrades to WebSocket.
// Query param: ?categories=uuid1,uuid2,... (expert passes their category IDs)
func (h *WebSocketHandler) HandleConnection(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// Parse optional category filter from query string
	var catIDs []uuid.UUID
	if raw := c.QueryParam("categories"); raw != "" {
		for _, part := range splitComma(raw) {
			if id, err := uuid.Parse(part); err == nil {
				catIDs = append(catIDs, id)
			}
		}
	}

	client := &Client{
		hub:         h.hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		categoryIDs: catIDs,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	return nil
}

func splitComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			p := s[start:i]
			if p != "" {
				parts = append(parts, p)
			}
			start = i + 1
		}
	}
	return parts
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
