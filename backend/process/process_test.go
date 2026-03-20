package process_test

import (
	"os"
	"testing"

	"github.com/user/server-manager/database"
	"github.com/user/server-manager/process"
)

func TestMain(m *testing.M) {
	// Initialize in-memory database
	database.InitDB(":memory:")
	code := m.Run()
	os.Exit(code)
}

func TestManagerStatusAndStop(t *testing.T) {
	mgr := process.GetManager()

	if mgr.Status() != process.StatusStopped {
		t.Errorf("expected initial status %s, got %s", process.StatusStopped, mgr.Status())
	}

	// Sending command while stopped should fail
	err := mgr.Command("help")
	if err == nil {
		t.Error("expected error sending command when stopped, got nil")
	}

	// Stopping while stopped should fail
	err = mgr.Stop()
	if err == nil {
		t.Error("expected error stopping when already stopped, got nil")
	}
}
