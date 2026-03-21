package minecraft

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/user/server-manager/internal/config"
	"github.com/user/server-manager/internal/service/state"
)

// FileService はマインクラフトサーバーに関する各種ファイル（設定ファイル、ホワイトリストなど）を操作するサービスを提供します。
type FileService struct{}

// NewFileService は FileService の新しいインスタンスを作成します。
func NewFileService() *FileService {
	return &FileService{}
}

// GetServerProperties は server.properties ファイルの内容を読み込み文字列として返します。
// エラーの場合は空文字列を返します。
func (fs *FileService) GetServerProperties() (string, error) {
	props, err := os.ReadFile(filepath.Join(config.GetServerDir(), "server.properties"))
	if err != nil {
		return "", err
	}
	return string(props), nil
}

// readJSONFile は指定されたJSONファイルを読み込み、インターフェース配列の形式にパースして返します。
// 例外的にファイルが存在しない場合は、パースエラーではなく空の配列を返します。
func (fs *FileService) readJSONFile(filename string) (interface{}, error) {
	data, err := os.ReadFile(filepath.Join(config.GetServerDir(), filename))
	if err != nil {
		return []interface{}{}, nil // ファイルが存在しない場合は空リストを返す
	}
	var result interface{}
	// 指定されたバイト列を任意の構造体にデコードする
	if err := json.Unmarshal(data, &result); err != nil {
		return []interface{}{}, err
	}
	return result, nil
}

// GetWhitelist は whitelist.json を読み込みインターフェース型で返します。
func (fs *FileService) GetWhitelist() (interface{}, error) {
	return fs.readJSONFile("whitelist.json")
}

// GetOps は ops.json (オペレーター情報) を読み込みインターフェース型で返します。
func (fs *FileService) GetOps() (interface{}, error) {
	return fs.readJSONFile("ops.json")
}

// GetBannedPlayers は banned-players.json を読み込みインターフェース型で返します。
func (fs *FileService) GetBannedPlayers() (interface{}, error) {
	return fs.readJSONFile("banned-players.json")
}

// GetOnlinePlayers は現在のオンライン状態にあるプレイヤーのリストを返します。
// nilスライスが返るのを避けるため、空の配列で初期化します。
func (fs *FileService) GetOnlinePlayers() []string {
	players := state.GetOnlinePlayers()
	if players == nil {
		return []string{} // null ではなく空配列を返す
	}
	return players
}

// GetServerLogs は latest.log ファイルの内容を読み込み、直近の200行分を配列として返します。
// ファイルが存在しない場合は空配列を返します。
func (fs *FileService) GetServerLogs() ([]string, error) {
	logFile := filepath.Join(config.GetServerDir(), "logs", "latest.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		// ログファイルがまだ存在しない場合は単に空で返す
		return []string{}, nil
	}

	// 大きなファイルの場合は本来tailを使うべきだが、簡素化のために全て読み込んでから分割する
	// （latest.log はログローテーションされるため通常は巨大になりすぎない前提）
	lines := []string{}
	currentLine := ""
	for _, b := range content {
		currentLine += string(b)
		if b == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// 巨大なペイロード（通信量）を回避するため、直近の200行に制限する
	if len(lines) > 200 {
		lines = lines[len(lines)-200:]
	}

	return lines, nil
}
