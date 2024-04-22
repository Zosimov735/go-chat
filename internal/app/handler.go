package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yourusername/yourproject/internal/db"
	"github.com/yourusername/yourproject/pkg/websocket"
	"log"
	"net/http"
)

type Handler struct {
	db         *db.Queries // Using sqlc or a similar library might be beneficial
	wsUpgrader websocket.Upgrader
}

func NewHandler(db *db.Queries) *Handler {
	return &Handler{
		db: db,
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // Configure appropriately
		},
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/ws", h.handleWebSocket)
}

// handleWebSocket upgrades the HTTP server connection to the WebSocket protocol.
func (h *Handler) handleWebSocket(c *gin.Context) {
	conn, err := h.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}
	// Initialize a new client session here
	client := websocket.NewClient(h.db, conn)
	go client.WritePump()
	go client.ReadPump()
}
