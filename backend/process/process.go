package process

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/user/server-manager/database"
)

type ServerStatus string

const (
	StatusStopped  ServerStatus = "Stopped"
	StatusStarting ServerStatus = "Starting"
	StatusRunning  ServerStatus = "Running"
)

type Manager struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	status ServerStatus
	mu     sync.Mutex
}

var instance *Manager

func GetManager() *Manager {
	if instance == nil {
		instance = &Manager{
			status: StatusStopped,
		}
	}
	return instance
}

func (m *Manager) Status() ServerStatus {
	m.mu.Lock()
	defer m.mu.Lock()
	return m.status
}

func (m *Manager) setStatus(s ServerStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = s
}

func (m *Manager) Start() error {
	m.mu.Lock()
	if m.status != StatusStopped {
		m.mu.Unlock()
		return errors.New("server is already running or starting")
	}
	m.status = StatusStarting
	m.mu.Unlock()

	// Ensure server directory exists
	_ = os.MkdirAll("../server", 0755)

	xms := database.GetSetting("java_xms")
	xmx := database.GetSetting("java_xmx")
	additionalArgs := database.GetSetting("java_args")

	args := []string{}
	if xms != "" {
		args = append(args, fmt.Sprintf("-Xms%s", xms))
	}
	if xmx != "" {
		args = append(args, fmt.Sprintf("-Xmx%s", xmx))
	}
	if additionalArgs != "" {
		args = append(args, strings.Split(additionalArgs, " ")...)
	}
	args = append(args, "-jar", "server.jar", "nogui")

	m.cmd = exec.Command("java", args...)
	m.cmd.Dir = "../server"

	stdin, err := m.cmd.StdinPipe()
	if err != nil {
		m.setStatus(StatusStopped)
		return err
	}
	m.stdin = stdin

	// For stdout, we'll let Minecraft write to logs/latest.log and tail it separately
	// or we can just pipe stdout to /dev/null if we only rely on the log file.
	// For now, let's keep stdout/stderr going to standard out so backend console sees it.
	m.cmd.Stdout = os.Stdout
	m.cmd.Stderr = os.Stderr

	if err := m.cmd.Start(); err != nil {
		m.setStatus(StatusStopped)
		return err
	}

	// Wait in a goroutine
	go func() {
		// Consider it 'Running' right after starting.
		// Realistically we'd parse logs for "Done" to be "Running", but this is simpler.
		time.Sleep(2 * time.Second)
		m.setStatus(StatusRunning)

		err := m.cmd.Wait()
		log.Printf("Server process exited: %v", err)

		m.setStatus(StatusStopped)
		m.cmd = nil
		m.stdin = nil
	}()

	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == StatusStopped || m.stdin == nil {
		return errors.New("server is not running")
	}

	// Send stop command
	_, err := io.WriteString(m.stdin, "stop\n")
	return err
}

func (m *Manager) Command(cmd string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status != StatusRunning || m.stdin == nil {
		return errors.New("server is not running")
	}

	_, err := io.WriteString(m.stdin, strings.TrimSpace(cmd)+"\n")
	return err
}
