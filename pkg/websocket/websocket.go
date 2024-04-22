package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/yourusername/yourproject/internal/db"
	"log"
	"time"
)

type Client struct {
	db   *db.Queries
	conn *websocket.Conn
	send chan []byte // Channel for outgoing messages
}

func NewClient(db *db.Queries, conn *websocket.Conn) *Client {
	return &Client{
		db:   db,
		conn: conn,
		send: make(chan []byte, 256), // Buffered channel
	}
}

// ReadPump listens for new messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Handle incoming messages, possibly store them in DB or broadcast
		log.Printf("recv: %s", message)
	}
}

// WritePump sends messages from the send channel to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The channel has been closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			// Keepalive pings can be sent here
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
