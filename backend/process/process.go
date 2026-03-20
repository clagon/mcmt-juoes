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

	"bufio"

	"github.com/user/server-manager/config"
	"github.com/user/server-manager/database"
	"github.com/user/server-manager/state"
	"github.com/user/server-manager/ws"
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
	defer m.mu.Unlock()
	return m.status
}

func (m *Manager) setStatus(s ServerStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = s
	state.SetServerStatus(string(s))
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
	serverDir := config.GetServerDir()
	_ = os.MkdirAll(serverDir, 0755)

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
	m.cmd.Dir = serverDir

	stdin, err := m.cmd.StdinPipe()
	if err != nil {
		m.setStatus(StatusStopped)
		return err
	}
	m.stdin = stdin

	stdoutPipe, err := m.cmd.StdoutPipe()
	if err != nil {
		m.setStatus(StatusStopped)
		return err
	}

	stderrPipe, err := m.cmd.StderrPipe()
	if err != nil {
		m.setStatus(StatusStopped)
		return err
	}

	if err := m.cmd.Start(); err != nil {
		m.setStatus(StatusStopped)
		return err
	}

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Log to backend console as well
			state.ParseLogForPlayers(line)
			ws.Broadcast("log", line)
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			ws.Broadcast("log", line)
		}
	}()

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
