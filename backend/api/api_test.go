package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/api"
)

func setupTestServerDir(t *testing.T) string {
	dir := t.TempDir()
	os.Setenv("SERVER_DIR", dir)
	return dir
}

func TestGetServerProperties(t *testing.T) {
	e := echo.New()
	dir := setupTestServerDir(t)
	defer os.Unsetenv("SERVER_DIR")

	// Create a dummy properties file
	propContent := "server-port=25565\nmotd=A Minecraft Server"
	err := os.WriteFile(filepath.Join(dir, "server.properties"), []byte(propContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/server/properties", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := api.GetServerProperties(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["content"] != propContent {
		t.Errorf("expected '%s', got '%s'", propContent, resp["content"])
	}
}

func TestGetServerProperties_Missing(t *testing.T) {
	e := echo.New()
	setupTestServerDir(t)
	defer os.Unsetenv("SERVER_DIR")

	// Don't create the file, verify it handles gracefully
	req := httptest.NewRequest(http.MethodGet, "/api/server/properties", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := api.GetServerProperties(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	var resp map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	if resp["content"] != "" {
		t.Errorf("expected empty string for missing file, got '%s'", resp["content"])
	}
}

func TestGetWhitelist(t *testing.T) {
	e := echo.New()
	dir := setupTestServerDir(t)
	defer os.Unsetenv("SERVER_DIR")

	whitelistContent := `[{"uuid":"123","name":"PlayerOne"}]`
	err := os.WriteFile(filepath.Join(dir, "whitelist.json"), []byte(whitelistContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/server/whitelist", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := api.GetWhitelist(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp []map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) != 1 || resp[0]["name"] != "PlayerOne" {
		t.Errorf("unexpected json response: %v", resp)
	}
}

func TestGetWhitelist_Missing(t *testing.T) {
	e := echo.New()
	setupTestServerDir(t)
	defer os.Unsetenv("SERVER_DIR")

	req := httptest.NewRequest(http.MethodGet, "/api/server/whitelist", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := api.GetWhitelist(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	var resp []interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	if len(resp) != 0 {
		t.Errorf("expected empty array for missing file, got len %d", len(resp))
	}
}
