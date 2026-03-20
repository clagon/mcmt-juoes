package config_test

import (
	"os"
	"testing"

	"github.com/user/server-manager/config"
)

func TestGetServerDir(t *testing.T) {
	// Test Default
	os.Unsetenv("SERVER_DIR")
	if dir := config.GetServerDir(); dir != "../server" {
		t.Errorf("expected default server dir '../server', got '%s'", dir)
	}

	// Test Env Var
	os.Setenv("SERVER_DIR", "/app/my-server")
	defer os.Unsetenv("SERVER_DIR")

	if dir := config.GetServerDir(); dir != "/app/my-server" {
		t.Errorf("expected env server dir '/app/my-server', got '%s'", dir)
	}
}

func TestGetDataDir(t *testing.T) {
	// Test Default
	os.Unsetenv("DATA_DIR")
	if dir := config.GetDataDir(); dir != "../data" {
		t.Errorf("expected default data dir '../data', got '%s'", dir)
	}

	// Test Env Var
	os.Setenv("DATA_DIR", "/app/my-data")
	defer os.Unsetenv("DATA_DIR")

	if dir := config.GetDataDir(); dir != "/app/my-data" {
		t.Errorf("expected env data dir '/app/my-data', got '%s'", dir)
	}
}
