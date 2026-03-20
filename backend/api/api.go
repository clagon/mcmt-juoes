package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/process"
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
	props, err := os.ReadFile("../server/server.properties")
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"content": ""})
	}
	return c.JSON(http.StatusOK, map[string]string{"content": string(props)})
}

// Player file parser helpers

func readJSONFile(filename string) (interface{}, error) {
	data, err := os.ReadFile("../server/" + filename)
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
