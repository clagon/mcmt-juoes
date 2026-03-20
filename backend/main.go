package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/user/server-manager/api"
	"github.com/user/server-manager/database"
	"github.com/user/server-manager/settings"
	"github.com/user/server-manager/ws"
)

func main() {
	// Initialize Database
	database.InitDB("../data/mcmt.db")

	// Start Background tasks
	go ws.StartStatusBroadcaster()
	go ws.StartLogTailer()

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

	// Websocket
	e.GET("/ws", ws.Handler)

	log.Println("Backend starting on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
