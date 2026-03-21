package state

import (
	"regexp"
	"sync"
)

var (
	onlinePlayers   = make(map[string]bool)
	onlinePlayersMu sync.Mutex

	joinRegex = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: (.+) joined the game`)
	leftRegex = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: (.+) left the game`)

	serverStatus   = "Stopped"
	serverStatusMu sync.Mutex
)

func GetServerStatus() string {
	serverStatusMu.Lock()
	defer serverStatusMu.Unlock()
	return serverStatus
}

func SetServerStatus(status string) {
	serverStatusMu.Lock()
	serverStatus = status
	serverStatusMu.Unlock()
}

func GetOnlinePlayers() []string {
	onlinePlayersMu.Lock()
	defer onlinePlayersMu.Unlock()

	var players []string
	for p := range onlinePlayers {
		players = append(players, p)
	}
	return players
}

func ParseLogForPlayers(line string) {
	if matches := joinRegex.FindStringSubmatch(line); len(matches) > 1 {
		player := matches[1]
		onlinePlayersMu.Lock()
		onlinePlayers[player] = true
		onlinePlayersMu.Unlock()
	} else if matches := leftRegex.FindStringSubmatch(line); len(matches) > 1 {
		player := matches[1]
		onlinePlayersMu.Lock()
		delete(onlinePlayers, player)
		onlinePlayersMu.Unlock()
	}
}

func ClearOnlinePlayers() {
	onlinePlayersMu.Lock()
	onlinePlayers = make(map[string]bool)
	onlinePlayersMu.Unlock()
}

func InitOnlinePlayersFromLog(lines []string) {
	ClearOnlinePlayers()
	for _, line := range lines {
		ParseLogForPlayers(line)
	}
}
