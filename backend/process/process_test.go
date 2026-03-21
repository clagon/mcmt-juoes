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

func TestManagerOperations(t *testing.T) {
	mgr := process.GetManager()

	tests := []struct {
		name      string
		operation func() error
		wantErr   bool
	}{
		{
			name: "Status is Stopped initially",
			operation: func() error {
				if mgr.Status() != process.StatusStopped {
					return os.ErrInvalid
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "Command fails when Stopped",
			operation: func() error {
				return mgr.Command("help")
			},
			wantErr: true,
		},
		{
			name: "Stop fails when Stopped",
			operation: func() error {
				return mgr.Stop()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.operation(); (err != nil) != tt.wantErr {
				t.Errorf("%s operation error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
