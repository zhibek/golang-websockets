package chat

import (
	"fmt"
)

type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection
}

func NewHub() Hub {
	return Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
}

func (h *Hub) Register(c *Connection) {
	h.register <- c
}

func (h *Hub) Unregister(c *Connection) {
	h.unregister <- c
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("Connect")
			h.connections[c] = true
		case c := <-h.unregister:
			fmt.Println("Disconnect")
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			fmt.Printf("Broadcasting: %s\n", m)
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}
