package ws

import (
	"bufio"
	"net/http"
	"os"
	"sync"
	"time"

	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/config"
	"github.com/user/server-manager/process"
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
	logCache  []string
	cacheMu   sync.Mutex
)

const maxLogCache = 100

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

func cacheLog(line string) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	logCache = append(logCache, line)
	if len(logCache) > maxLogCache {
		logCache = logCache[1:]
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
		Data: process.GetManager().Status(),
	})

	// Send cached logs
	cacheMu.Lock()
	for _, line := range logCache {
		_ = ws.WriteJSON(Message{Type: "log", Data: line})
	}
	cacheMu.Unlock()

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

// StartStatusBroadcaster polls server status and broadcasts it if changed.
func StartStatusBroadcaster() {
	lastStatus := process.GetManager().Status()
	for {
		currentStatus := process.GetManager().Status()
		if currentStatus != lastStatus {
			Broadcast("status", currentStatus)
			lastStatus = currentStatus
		}
		time.Sleep(1 * time.Second)
	}
}

func StartLogTailer() {
	logFile := filepath.Join(config.GetServerDir(), "logs", "latest.log")

	var file *os.File
	var err error

	for {
		file, err = os.Open(logFile)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	defer file.Close()

	// Read from the beginning on start
	_, _ = file.Seek(0, 0)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(500 * time.Millisecond) // wait for more data

			// check if file was rotated/truncated
			info, statErr := file.Stat()
			if statErr == nil && info.Size() < getFileSize(file) {
				_, _ = file.Seek(0, 0)
				reader.Reset(file)
			}
			continue
		}

		cacheLog(line)
		Broadcast("log", line)
	}
}

func getFileSize(f *os.File) int64 {
	offset, _ := f.Seek(0, 1)
	return offset
}
