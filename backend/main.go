package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/user/server-manager/internal/config"
	"github.com/user/server-manager/internal/handler/api"
	"github.com/user/server-manager/internal/handler/settings"
	"github.com/user/server-manager/internal/repository"
	"github.com/user/server-manager/internal/service/state"
	"github.com/user/server-manager/internal/service/ws"
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
	repository.InitDB(filepath.Join(config.GetDataDir(), "mcmt.db"))

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

	// Handlers
	mcHandler := api.NewMinecraftHandler()

	// Settings API
	e.GET("/api/settings", settings.GetSettings)
	e.PUT("/api/settings", settings.UpdateSettings)

	// Server API
	e.POST("/api/server/start", mcHandler.StartServer)
	e.POST("/api/server/stop", mcHandler.StopServer)
	e.POST("/api/server/command", mcHandler.CommandServer)

	// File API
	e.GET("/api/server/properties", mcHandler.GetServerProperties)
	e.GET("/api/server/whitelist", mcHandler.GetWhitelist)
	e.GET("/api/server/ops", mcHandler.GetOps)
	e.GET("/api/server/banned-players", mcHandler.GetBannedPlayers)
	e.GET("/api/server/online", mcHandler.GetOnlinePlayers)
	e.GET("/api/server/logs", mcHandler.GetServerLogs)

	// Websocket
	e.GET("/ws", ws.Handler)

	log.Println("Backend starting on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
