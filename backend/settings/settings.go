package settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/database"
)

type Settings struct {
	JavaXms  string `json:"java_xms"`
	JavaXmx  string `json:"java_xmx"`
	JavaArgs string `json:"java_args"`
}

func GetSettings(c echo.Context) error {
	s := Settings{
		JavaXms:  database.GetSetting("java_xms"),
		JavaXmx:  database.GetSetting("java_xmx"),
		JavaArgs: database.GetSetting("java_args"),
	}
	return c.JSON(http.StatusOK, s)
}

func UpdateSettings(c echo.Context) error {
	s := new(Settings)
	if err := c.Bind(s); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	_ = database.SetSetting("java_xms", s.JavaXms)
	_ = database.SetSetting("java_xmx", s.JavaXmx)
	_ = database.SetSetting("java_args", s.JavaArgs)

	return c.JSON(http.StatusOK, s)
}
