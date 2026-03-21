package minecraft

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/user/server-manager/internal/config"
	"github.com/user/server-manager/internal/repository"
	"github.com/user/server-manager/internal/service/state"
	"github.com/user/server-manager/internal/service/ws"
)

// ServerStatus はマインクラフトサーバーの現在の状態を表す型です。
type ServerStatus string

const (
	// StatusStopped はサーバーが停止している状態を示します。
	StatusStopped ServerStatus = "Stopped"
	// StatusStarting はサーバーが起動中の状態を示します。
	StatusStarting ServerStatus = "Starting"
	// StatusRunning はサーバーが稼働中の状態を示します。
	StatusRunning ServerStatus = "Running"
)

// ProcessManager はマインクラフトサーバーのプロセス（起動、停止、コマンド実行）を管理します。
type ProcessManager struct {
	cmd    *exec.Cmd      // 実行中のJavaプロセス
	stdin  io.WriteCloser // プロセスへの標準入力（コマンド送信用）
	status ServerStatus   // 現在のサーバーステータス
	mu     sync.Mutex     // 状態変更時の排他制御用ミューテックス
}

// instance は ProcessManager のシングルトンインスタンスです。
var instance *ProcessManager

// GetManager は ProcessManager のシングルトンインスタンスを返します。
// アプリケーション全体で一つのプロセスを共有・管理するために使用されます。
func GetManager() *ProcessManager {
	if instance == nil {
		instance = &ProcessManager{
			status: StatusStopped,
		}
	}
	return instance
}

// Status は現在のマインクラフトサーバーの状態（Stopped, Starting, Running）を返します。
// スレッドセーフに状態を取得できます。
func (m *ProcessManager) Status() ServerStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

// setStatus はサーバーステータスを更新し、内部状態の同期とグローバルな状態管理への反映を行います。
func (m *ProcessManager) setStatus(s ServerStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = s
	state.SetServerStatus(string(s))
}

// Start はマインクラフトサーバーのJavaプロセスを起動します。
// 既に起動中・起動処理中の場合はエラーを返します。
func (m *ProcessManager) Start() error {
	m.mu.Lock()
	if m.status != StatusStopped {
		m.mu.Unlock()
		return errors.New("server is already running or starting")
	}
	m.status = StatusStarting
	m.mu.Unlock()

	// サーバー実行用のディレクトリの存在を保証する
	serverDir := config.GetServerDir()
	_ = os.MkdirAll(serverDir, 0755)

	// DBからJavaの起動オプションを取得する
	xms := repository.GetSetting("java_xms")
	xmx := repository.GetSetting("java_xmx")
	additionalArgs := repository.GetSetting("java_args")

	args := []string{}
	if xms != "" {
		args = append(args, fmt.Sprintf("-Xms%s", xms))
	}
	if xmx != "" {
		args = append(args, fmt.Sprintf("-Xmx%s", xmx))
	}
	if additionalArgs != "" {
		// 追加の引数がある場合はスペースで分割して追加
		args = append(args, strings.Split(additionalArgs, " ")...)
	}
	// サーバーjarファイルとnoguiオプション（GUI不要）を指定
	args = append(args, "-jar", "server.jar", "nogui")

	// javaコマンドの実行準備
	m.cmd = exec.Command("java", args...)
	m.cmd.Dir = serverDir

	// 標準入出力をパイプで繋ぐ（コマンド送信とログ取得用）
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

	// プロセスの非同期起動
	if err := m.cmd.Start(); err != nil {
		m.setStatus(StatusStopped)
		return err
	}

	// 標準出力（標準ログ）をストリームして読み取るゴルーチン
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)              // バックエンドのコンソールにも出力
			state.ParseLogForPlayers(line) // ログイン・ログアウト判定
			ws.Broadcast("log", line)      // フロントエンドへログを配信
		}
	}()

	// 標準エラー出力（エラーログ）をストリームして読み取るゴルーチン
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			ws.Broadcast("log", line)
		}
	}()

	// プロセスの終了を待ち受けるゴルーチン
	go func() {
		// 起動直後にRunning状態とする（本来はログから"Done"を判定する方が正確だが、簡易化のため）
		time.Sleep(2 * time.Second)
		m.setStatus(StatusRunning)

		// プロセスの終了を待機（stopコマンド等により終了するまでブロック）
		err := m.cmd.Wait()
		log.Printf("Server process exited: %v", err)

		// プロセス終了後のクリーンアップ処理
		m.setStatus(StatusStopped)
		m.cmd = nil
		m.stdin = nil
	}()

	return nil
}

// Stop は稼働中のマインクラフトサーバーを安全に停止します。
// 標準入力を通じて "stop" コマンドを送信します。
func (m *ProcessManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == StatusStopped || m.stdin == nil {
		return errors.New("server is not running")
	}

	// "stop" コマンドを送信して安全に終了させる
	_, err := io.WriteString(m.stdin, "stop\n")
	return err
}

// Command は稼働中のマインクラフトサーバーに任意のコマンドを送信します。
func (m *ProcessManager) Command(cmd string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status != StatusRunning || m.stdin == nil {
		return errors.New("server is not running")
	}

	// コマンド文字列の末尾に改行を付与して標準入力に書き込む
	_, err := io.WriteString(m.stdin, strings.TrimSpace(cmd)+"\n")
	return err
}
