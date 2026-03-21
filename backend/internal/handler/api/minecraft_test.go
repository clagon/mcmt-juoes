package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/internal/config"
	"github.com/user/server-manager/internal/handler/api"
)

func TestGetServerProperties(t *testing.T) {
	// Setup test directory
	tmpDir, _ := os.MkdirTemp("", "mcmt-test")
	defer os.RemoveAll(tmpDir)

	os.Setenv("SERVER_DIR", tmpDir)
	defer os.Unsetenv("SERVER_DIR")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/server/properties", nil)

	h := api.NewMinecraftHandler()

	t.Run("File Exists", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		content := "motd=A Minecraft Server\nserver-port=25565"

		// 確実に環境変数が読み込まれるようにconfigの初期化等が行われているか確認が必要だが、GetServerDirは毎度Getenvを読んでいるため問題ない
		os.WriteFile(filepath.Join(config.GetServerDir(), "server.properties"), []byte(content), 0644)

		if err := h.GetServerProperties(c); err != nil {
			t.Errorf("GetServerProperties returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		expectedBody := `{"content":"motd=A Minecraft Server\nserver-port=25565"}` + "\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("File Missing", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		os.Remove(filepath.Join(config.GetServerDir(), "server.properties"))

		if err := h.GetServerProperties(c); err != nil {
			t.Errorf("GetServerProperties returned error: %v", err)
		}

		expectedBody := `{"content":""}` + "\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}

func TestGetWhitelist(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "mcmt-test")
	defer os.RemoveAll(tmpDir)
	os.Setenv("SERVER_DIR", tmpDir)
	defer os.Unsetenv("SERVER_DIR")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/server/whitelist", nil)

	h := api.NewMinecraftHandler()

	t.Run("File Exists and Valid", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		content := `[{"uuid":"123","name":"player1"}]`
		os.WriteFile(filepath.Join(config.GetServerDir(), "whitelist.json"), []byte(content), 0644)

		if err := h.GetWhitelist(c); err != nil {
			t.Errorf("GetWhitelist returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		expectedBody := `[{"name":"player1","uuid":"123"}]` + "\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("File Missing", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		os.Remove(filepath.Join(config.GetServerDir(), "whitelist.json"))

		if err := h.GetWhitelist(c); err != nil {
			t.Errorf("GetWhitelist returned error: %v", err)
		}

		expectedBody := `[]` + "\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}
