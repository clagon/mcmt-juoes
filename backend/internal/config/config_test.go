package config_test

import (
	"os"
	"testing"

	"github.com/user/server-manager/internal/config"
)

func TestGetServerDir(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		setEnv bool
		want   string
	}{
		{
			name:   "Default",
			setEnv: false,
			want:   "../server",
		},
		{
			name:   "Custom Env",
			envVal: "/app/custom-server",
			setEnv: true,
			want:   "/app/custom-server",
		},
		{
			name:   "Empty Env (fallback to default)",
			envVal: "",
			setEnv: true,
			want:   "../server",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv("SERVER_DIR", tt.envVal)
			} else {
				os.Unsetenv("SERVER_DIR")
			}
			defer os.Unsetenv("SERVER_DIR")

			if got := config.GetServerDir(); got != tt.want {
				t.Errorf("GetServerDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDataDir(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		setEnv bool
		want   string
	}{
		{
			name:   "Default",
			setEnv: false,
			want:   "../data",
		},
		{
			name:   "Custom Env",
			envVal: "/app/custom-data",
			setEnv: true,
			want:   "/app/custom-data",
		},
		{
			name:   "Empty Env (fallback to default)",
			envVal: "",
			setEnv: true,
			want:   "../data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv("DATA_DIR", tt.envVal)
			} else {
				os.Unsetenv("DATA_DIR")
			}
			defer os.Unsetenv("DATA_DIR")

			if got := config.GetDataDir(); got != tt.want {
				t.Errorf("GetDataDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
