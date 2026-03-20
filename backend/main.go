package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/user/server-manager/api"
	"github.com/user/server-manager/config"
	"github.com/user/server-manager/database"
	"github.com/user/server-manager/settings"
	"github.com/user/server-manager/state"
	"github.com/user/server-manager/ws"
)

func StartStatusBroadcaster() {
	lastStatus := state.GetServerStatus()
	for {
		currentStatus := state.GetServerStatus()
		if currentStatus != lastStatus {
			ws.Broadcast("status", currentStatus)

			if currentStatus == "Stopped" {
				state.ClearOnlinePlayers()
			}

			lastStatus = currentStatus
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Initialize Database
	database.InitDB(filepath.Join(config.GetDataDir(), "mcmt.db"))

	// Init online players from existing latest.log if server might be running
	logFile := filepath.Join(config.GetServerDir(), "logs", "latest.log")
	content, _ := os.ReadFile(logFile)
	lines := []string{}
	currentLine := ""
	for _, b := range content {
		currentLine += string(b)
		if b == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	state.InitOnlinePlayersFromLog(lines)

	// Start Background tasks
	go StartStatusBroadcaster()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Settings API
	e.GET("/api/settings", settings.GetSettings)
	e.PUT("/api/settings", settings.UpdateSettings)

	// Server API
	e.POST("/api/server/start", api.StartServer)
	e.POST("/api/server/stop", api.StopServer)
	e.POST("/api/server/command", api.CommandServer)

	// File API
	e.GET("/api/server/properties", api.GetServerProperties)
	e.GET("/api/server/whitelist", api.GetWhitelist)
	e.GET("/api/server/ops", api.GetOps)
	e.GET("/api/server/banned-players", api.GetBannedPlayers)
	e.GET("/api/server/online", api.GetOnlinePlayers)
	e.GET("/api/server/logs", api.GetServerLogs)

	// Websocket
	e.GET("/ws", ws.Handler)

	log.Println("Backend starting on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
