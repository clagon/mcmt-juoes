package state

import (
	"regexp"
	"sync"
)

var (
	// onlinePlayers は現在オンラインであるプレイヤーを管理するマップです
	onlinePlayers   = make(map[string]bool)
	onlinePlayersMu sync.Mutex

	// joinRegex はプレイヤーのログインを検出する正規表現（サーバーのログ形式に依存）
	joinRegex = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: (.+) joined the game`)
	// leftRegex はプレイヤーのログアウトを検出する正規表現
	leftRegex = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread\/INFO\]: (.+) left the game`)

	// serverStatus はマインクラフトサーバーの現在のステータスを保持します（例: Stopped, Starting, Running）
	serverStatus   = "Stopped"
	serverStatusMu sync.Mutex
)

// GetServerStatus は現在のサーバーステータスをスレッドセーフに取得します。
func GetServerStatus() string {
	serverStatusMu.Lock()
	defer serverStatusMu.Unlock()
	return serverStatus
}

// SetServerStatus は現在のサーバーステータスを更新します。プロセス管理クラスから呼ばれます。
func SetServerStatus(status string) {
	serverStatusMu.Lock()
	serverStatus = status
	serverStatusMu.Unlock()
}

// GetOnlinePlayers は現在オンライン状態にある全プレイヤーの名前をスレッドセーフに取得し、スライスで返します。
func GetOnlinePlayers() []string {
	onlinePlayersMu.Lock()
	defer onlinePlayersMu.Unlock()

	var players []string
	for p := range onlinePlayers {
		players = append(players, p)
	}
	return players
}

// ParseLogForPlayers は標準出力されるログの各行を解析し、プレイヤーのログイン・ログアウト状態を更新します。
func ParseLogForPlayers(line string) {
	// マッチした場合は第1キャプチャグループ（[1]）からプレイヤー名を抽出
	if matches := joinRegex.FindStringSubmatch(line); len(matches) > 1 {
		player := matches[1]
		onlinePlayersMu.Lock()
		onlinePlayers[player] = true // ログイン状態にする
		onlinePlayersMu.Unlock()
	} else if matches := leftRegex.FindStringSubmatch(line); len(matches) > 1 {
		player := matches[1]
		onlinePlayersMu.Lock()
		delete(onlinePlayers, player) // ログアウト状態としてマップから削除
		onlinePlayersMu.Unlock()
	}
}

// ClearOnlinePlayers はオンラインプレイヤー情報をリセットして初期化します。サーバー停止時などに呼ばれます。
func ClearOnlinePlayers() {
	onlinePlayersMu.Lock()
	onlinePlayers = make(map[string]bool)
	onlinePlayersMu.Unlock()
}

// InitOnlinePlayersFromLog はログファイルの内容を解析してオンラインプレイヤー状態を初期化します。
// （サーバー起動中にバックエンドアプリが再起動された場合などを想定しています）
func InitOnlinePlayersFromLog(lines []string) {
	ClearOnlinePlayers()
	for _, line := range lines {
		ParseLogForPlayers(line)
	}
}
