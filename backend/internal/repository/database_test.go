package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDBAndSettings(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcmt-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	InitDB(dbPath)

	// verify defaults
	xms := GetSetting("java_xms")
	if xms != "2G" {
		t.Errorf("Expected 2G, got %s", xms)
	}

	xmx := GetSetting("java_xmx")
	if xmx != "2G" {
		t.Errorf("Expected 2G, got %s", xmx)
	}

	args := GetSetting("java_args")
	if args != "" {
		t.Errorf("Expected empty args, got %s", args)
	}

	// Update setting
	err = SetSetting("java_xms", "4G")
	if err != nil {
		t.Errorf("Failed to set setting: %v", err)
	}

	xmsUpdated := GetSetting("java_xms")
	if xmsUpdated != "4G" {
		t.Errorf("Expected 4G, got %s", xmsUpdated)
	}
}
