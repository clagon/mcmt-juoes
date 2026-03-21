package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/internal/service/state"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local tool
	},
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func Broadcast(msgType string, data interface{}) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	msg := Message{Type: msgType, Data: data}

	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func Handler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	// Send current status on connect
	_ = ws.WriteJSON(Message{
		Type: "status",
		Data: state.GetServerStatus(),
	})

	for {
		// Read message loop (keep alive)
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	clientsMu.Lock()
	delete(clients, ws)
	clientsMu.Unlock()

	return nil
}
