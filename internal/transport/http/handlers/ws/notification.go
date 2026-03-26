package ws

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeTimeout   = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 4096
	writeChanSize  = 8
)

type incomingMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type outgoingMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Notification struct {
	upgrader websocket.Upgrader
}

func NewNotification(allowedOrigins []string) *Notification {
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = struct{}{}
	}

	return &Notification{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 10 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				_, ok := originSet[origin]
				return ok
			},
		},
	}
}

// Connect godoc
//
//	@Summary	WebSocket connection for real-time notifications
//	@Tags		ws
//	@Router		/ws/notifications [get]
func (h *Notification) Connect(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "ws: upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	ctx := c.Request.Context()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { return conn.SetReadDeadline(time.Now().Add(pongWait)) })

	writeCh := make(chan outgoingMessage, writeChanSize)
	done := make(chan struct{})

	go h.writeLoop(ctx, conn, writeCh, done)
	defer close(done)

	h.readLoop(ctx, conn, writeCh)
}

// readLoop reads incoming JSON messages and forwards replies to writeLoop via writeCh.
// It is the only goroutine that reads from conn.
func (h *Notification) readLoop(ctx context.Context, conn *websocket.Conn, writeCh chan<- outgoingMessage) {
	for {
		var msg incomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				slog.ErrorContext(ctx, "ws: read error", "error", err)
			}
			return
		}

		select {
		case writeCh <- outgoingMessage{Type: "echo", Payload: msg.Payload}:
		case <-ctx.Done():
			return
		}
	}
}

// writeLoop is the only goroutine that writes to conn, eliminating concurrent write races.
// It also sends periodic pings to keep the connection alive.
func (h *Notification) writeLoop(
	ctx context.Context,
	conn *websocket.Conn,
	writeCh <-chan outgoingMessage,
	done <-chan struct{},
) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg := <-writeCh:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteJSON(msg); err != nil {
				slog.ErrorContext(ctx, "ws: write error", "error", err)
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-done:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		}
	}
}
