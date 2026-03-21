package settings_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/internal/handler/settings"
	"github.com/user/server-manager/internal/repository"
)

func TestMain(m *testing.M) {
	repository.InitDB(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestGetSettings(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		wantStatus  int
		wantJavaXms string
		wantJavaXmx string
	}{
		{
			name: "Success Returns 1G 2G",
			setup: func() {
				_ = repository.SetSetting("java_xms", "1G")
				_ = repository.SetSetting("java_xmx", "2G")
				_ = repository.SetSetting("java_args", "-XX:+UseG1GC")
			},
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantJavaXms: "1G",
			wantJavaXmx: "2G",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			args := args{c: c}

			if err := settings.GetSettings(args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantStatus {
				t.Errorf("GetSettings() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			var s settings.Settings
			if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if s.JavaXms != tt.wantJavaXms {
				t.Errorf("GetSettings() JavaXms = %v, want %v", s.JavaXms, tt.wantJavaXms)
			}
			if s.JavaXmx != tt.wantJavaXmx {
				t.Errorf("GetSettings() JavaXmx = %v, want %v", s.JavaXmx, tt.wantJavaXmx)
			}
		})
	}
}

func TestUpdateSettings(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name       string
		reqBody    interface{}
		wantErr    bool
		wantStatus int
		wantDBXms  string
	}{
		{
			name: "Success Update Settings to 3G",
			reqBody: settings.Settings{
				JavaXms:  "3G",
				JavaXmx:  "6G",
				JavaArgs: "-Dfoo=bar",
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantDBXms:  "3G",
		},
		{
			name:       "Invalid JSON Body",
			reqBody:    "not-a-json",
			wantErr:    false, // Handler handles the error internally by returning a string status
			wantStatus: http.StatusBadRequest,
			wantDBXms:  "3G", // DB state shouldn't change from last successful update
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			bodyBytes, _ := json.Marshal(tt.reqBody)
			if s, ok := tt.reqBody.(string); ok {
				bodyBytes = []byte(s) // For invalid JSON case
			}

			req := httptest.NewRequest(http.MethodPut, "/api/settings", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			args := args{c: c}

			if err := settings.UpdateSettings(args.c); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantStatus {
				t.Errorf("UpdateSettings() status = %v, want %v", rec.Code, tt.wantStatus)
			}
			if dbVal := repository.GetSetting("java_xms"); dbVal != tt.wantDBXms {
				t.Errorf("UpdateSettings() DB java_xms = %v, want %v", dbVal, tt.wantDBXms)
			}
		})
	}
}
