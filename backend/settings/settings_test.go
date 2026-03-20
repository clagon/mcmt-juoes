package settings_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/database"
	"github.com/user/server-manager/settings"
)

func TestMain(m *testing.M) {
	// Setup test database
	database.InitDB(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestGetSettings(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set initial values
	database.SetSetting("java_xms", "1G")
	database.SetSetting("java_xmx", "2G")
	database.SetSetting("java_args", "-XX:+UseG1GC")

	if err := settings.GetSettings(c); err != nil {
		t.Fatalf("GetSettings handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var s settings.Settings
	if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if s.JavaXms != "1G" {
		t.Errorf("expected JavaXms=1G, got %s", s.JavaXms)
	}
	if s.JavaXmx != "2G" {
		t.Errorf("expected JavaXmx=2G, got %s", s.JavaXmx)
	}
}

func TestUpdateSettings(t *testing.T) {
	e := echo.New()

	newSettings := settings.Settings{
		JavaXms: "3G",
		JavaXmx: "6G",
		JavaArgs: "-Dfoo=bar",
	}

	body, _ := json.Marshal(newSettings)
	req := httptest.NewRequest(http.MethodPut, "/api/settings", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := settings.UpdateSettings(c); err != nil {
		t.Fatalf("UpdateSettings handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if database.GetSetting("java_xms") != "3G" {
		t.Errorf("setting did not update in DB")
	}
}
