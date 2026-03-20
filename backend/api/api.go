package api

import (
	"encoding/json"
	"net/http"
	"os"

	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/config"
	"github.com/user/server-manager/process"
	"github.com/user/server-manager/state"
)

func StartServer(c echo.Context) error {
	err := process.GetManager().Start()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Server started"})
}

func StopServer(c echo.Context) error {
	err := process.GetManager().Stop()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Server stopped"})
}

func CommandServer(c echo.Context) error {
	var body struct {
		Command string `json:"command"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	err := process.GetManager().Command(body.Command)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Command sent"})
}

func GetServerProperties(c echo.Context) error {
	props, err := os.ReadFile(filepath.Join(config.GetServerDir(), "server.properties"))
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"content": ""})
	}
	return c.JSON(http.StatusOK, map[string]string{"content": string(props)})
}

// Player file parser helpers

func readJSONFile(filename string) (interface{}, error) {
	data, err := os.ReadFile(filepath.Join(config.GetServerDir(), filename))
	if err != nil {
		return []interface{}{}, nil // Return empty list if file doesn't exist
	}
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return []interface{}{}, err
	}
	return result, nil
}

func GetWhitelist(c echo.Context) error {
	data, err := readJSONFile("whitelist.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func GetOps(c echo.Context) error {
	data, err := readJSONFile("ops.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func GetBannedPlayers(c echo.Context) error {
	data, err := readJSONFile("banned-players.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func GetOnlinePlayers(c echo.Context) error {
	players := state.GetOnlinePlayers()
	if players == nil {
		players = []string{} // Return empty array instead of null
	}
	return c.JSON(http.StatusOK, players)
}

func GetServerLogs(c echo.Context) error {
	logFile := filepath.Join(config.GetServerDir(), "logs", "latest.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		// Log file might not exist yet, just return empty
		return c.JSON(http.StatusOK, []string{})
	}

	// We'll return the lines. For large files we'd tail, but we just read all here for simplicity.
	// If it's too large, ideally we'd limit to last N lines.
	// For now, since 'latest.log' usually isn't massive due to log rotation,
	// returning split strings is fine for the terminal view.
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

	// Optionally limit to the last 200 lines to avoid massive payloads
	if len(lines) > 200 {
		lines = lines[len(lines)-200:]
	}

	return c.JSON(http.StatusOK, lines)
}
