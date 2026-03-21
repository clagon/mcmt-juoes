package state

import (
	"testing"
)

func TestServerStatus(t *testing.T) {
	// 初期状態は Stopped と定義されているが、テスト実行順に依存しないように再設定
	SetServerStatus("Stopped")

	if status := GetServerStatus(); status != "Stopped" {
		t.Errorf("Expected status Stopped, got %s", status)
	}

	SetServerStatus("Running")
	if status := GetServerStatus(); status != "Running" {
		t.Errorf("Expected status Running, got %s", status)
	}
}

func TestOnlinePlayers(t *testing.T) {
	ClearOnlinePlayers()

	players := GetOnlinePlayers()
	if len(players) != 0 {
		t.Errorf("Expected 0 players initially, got %d", len(players))
	}

	joinLine := "[12:34:56] [Server thread/INFO]: Player1 joined the game"
	ParseLogForPlayers(joinLine)

	players = GetOnlinePlayers()
	if len(players) != 1 || players[0] != "Player1" {
		t.Errorf("Expected [Player1], got %v", players)
	}

	joinLine2 := "[12:35:00] [Server thread/INFO]: Player2 joined the game"
	ParseLogForPlayers(joinLine2)

	players = GetOnlinePlayers()
	if len(players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(players))
	}
	// map keys are not ordered, just check length and existence
	found1 := false
	found2 := false
	for _, p := range players {
		if p == "Player1" {
			found1 = true
		}
		if p == "Player2" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("Expected Player1 and Player2, got %v", players)
	}

	leftLine := "[12:36:00] [Server thread/INFO]: Player1 left the game"
	ParseLogForPlayers(leftLine)

	players = GetOnlinePlayers()
	if len(players) != 1 || players[0] != "Player2" {
		t.Errorf("Expected [Player2], got %v", players)
	}

	ClearOnlinePlayers()
	if len(GetOnlinePlayers()) != 0 {
		t.Errorf("Expected 0 players after clear")
	}
}

func TestInitOnlinePlayersFromLog(t *testing.T) {
	lines := []string{
		"[12:34:56] [Server thread/INFO]: PlayerA joined the game",
		"[12:35:00] [Server thread/INFO]: PlayerB joined the game",
		"[12:36:00] [Server thread/INFO]: PlayerA left the game",
		"Some random log line",
	}

	InitOnlinePlayersFromLog(lines)

	players := GetOnlinePlayers()
	if len(players) != 1 || players[0] != "PlayerB" {
		t.Errorf("Expected [PlayerB], got %v", players)
	}
}
