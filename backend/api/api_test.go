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
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name        string
		createFile  bool
		fileContent string
		wantErr     bool
		wantStatus  int
		wantBody    map[string]string
	}{
		{
			name:        "File Exists",
			createFile:  true,
			fileContent: "server-port=25565\nmotd=A Minecraft Server",
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantBody:    map[string]string{"content": "server-port=25565\nmotd=A Minecraft Server"},
		},
		{
			name:       "File Missing",
			createFile: false,
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"content": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := setupTestServerDir(t)
			defer os.Unsetenv("SERVER_DIR")

			if tt.createFile {
				err := os.WriteFile(filepath.Join(dir, "server.properties"), []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/server/properties", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			args := args{c: c}

			if err := api.GetServerProperties(args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetServerProperties() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("GetServerProperties() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			var resp map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if resp["content"] != tt.wantBody["content"] {
				t.Errorf("GetServerProperties() body = %v, want %v", resp, tt.wantBody)
			}
		})
	}
}

func TestGetWhitelist(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name          string
		createFile    bool
		fileContent   string
		wantErr       bool
		wantStatus    int
		wantBodyLen   int
		wantBodyFirst map[string]string
	}{
		{
			name:          "File Exists and Valid",
			createFile:    true,
			fileContent:   `[{"uuid":"123","name":"PlayerOne"}]`,
			wantErr:       false,
			wantStatus:    http.StatusOK,
			wantBodyLen:   1,
			wantBodyFirst: map[string]string{"name": "PlayerOne", "uuid": "123"},
		},
		{
			name:        "File Missing",
			createFile:  false,
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantBodyLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := setupTestServerDir(t)
			defer os.Unsetenv("SERVER_DIR")

			if tt.createFile {
				err := os.WriteFile(filepath.Join(dir, "whitelist.json"), []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/server/whitelist", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			args := args{c: c}

			if err := api.GetWhitelist(args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetWhitelist() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("GetWhitelist() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			var resp []map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if len(resp) != tt.wantBodyLen {
				t.Errorf("GetWhitelist() array len = %v, want %v", len(resp), tt.wantBodyLen)
			}

			if tt.wantBodyLen > 0 {
				if resp[0]["name"] != tt.wantBodyFirst["name"] {
					t.Errorf("GetWhitelist() first item name = %v, want %v", resp[0]["name"], tt.wantBodyFirst["name"])
				}
			}
		})
	}
}
