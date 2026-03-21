package minecraft

import (
	"testing"
)

func TestManagerOperations(t *testing.T) {
	manager := GetManager()

	t.Run("Status is Stopped initially", func(t *testing.T) {
		if manager.Status() != StatusStopped {
			t.Errorf("Expected status Stopped, got %s", manager.Status())
		}
	})

	t.Run("Command fails when Stopped", func(t *testing.T) {
		err := manager.Command("help")
		if err == nil {
			t.Error("Expected error when sending command to stopped server")
		}
	})

	t.Run("Stop fails when Stopped", func(t *testing.T) {
		err := manager.Stop()
		if err == nil {
			t.Error("Expected error when stopping a stopped server")
		}
	})
}
