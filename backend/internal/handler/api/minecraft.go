package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user/server-manager/internal/service/minecraft"
)

// MinecraftHandler はマインクラフトサーバーに対するHTTPリクエストを処理するハンドラ群です。
type MinecraftHandler struct {
	fileService *minecraft.FileService // ファイル操作に関するビジネスロジックを委譲するサービス
}

// NewMinecraftHandler は MinecraftHandler の新しいインスタンスを作成して返します。
// （簡易的なDIとして内部でサービスを初期化しています）
func NewMinecraftHandler() *MinecraftHandler {
	return &MinecraftHandler{
		fileService: minecraft.NewFileService(),
	}
}

// StartServer はサーバー起動のエンドポイント（POST /api/server/start）の処理を行います。
func (h *MinecraftHandler) StartServer(c echo.Context) error {
	// シングルトンの ProcessManager を使用して起動
	err := minecraft.GetManager().Start()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Server started"})
}

// StopServer はサーバー停止のエンドポイント（POST /api/server/stop）の処理を行います。
func (h *MinecraftHandler) StopServer(c echo.Context) error {
	// シングルトンの ProcessManager を使用して安全に停止
	err := minecraft.GetManager().Stop()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Server stopped"})
}

// CommandServer は任意のコマンド送信エンドポイント（POST /api/server/command）の処理を行います。
func (h *MinecraftHandler) CommandServer(c echo.Context) error {
	// JSONリクエストボディの構造体定義
	var body struct {
		Command string `json:"command"`
	}
	// リクエストボディのバインディング（パース）
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	// 抽出したコマンドをサーバーの標準入力に送信する
	err := minecraft.GetManager().Command(body.Command)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Command sent"})
}

// GetServerProperties は server.properties の内容を取得して返すエンドポイント（GET /api/server/properties）です。
func (h *MinecraftHandler) GetServerProperties(c echo.Context) error {
	props, err := h.fileService.GetServerProperties()
	if err != nil {
		// 失敗した場合はエラーではなく空の内容として返す（ファイル未生成の可能性があるため）
		return c.JSON(http.StatusOK, map[string]string{"content": ""})
	}
	return c.JSON(http.StatusOK, map[string]string{"content": props})
}

// GetWhitelist はホワイトリスト（whitelist.json）を返すエンドポイント（GET /api/server/whitelist）です。
func (h *MinecraftHandler) GetWhitelist(c echo.Context) error {
	data, err := h.fileService.GetWhitelist()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

// GetOps はオペレーターリスト（ops.json）を返すエンドポイント（GET /api/server/ops）です。
func (h *MinecraftHandler) GetOps(c echo.Context) error {
	data, err := h.fileService.GetOps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

// GetBannedPlayers はBANリスト（banned-players.json）を返すエンドポイント（GET /api/server/banned-players）です。
func (h *MinecraftHandler) GetBannedPlayers(c echo.Context) error {
	data, err := h.fileService.GetBannedPlayers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

// GetOnlinePlayers は現在のオンラインプレイヤーリストを返すエンドポイント（GET /api/server/online）です。
func (h *MinecraftHandler) GetOnlinePlayers(c echo.Context) error {
	players := h.fileService.GetOnlinePlayers()
	return c.JSON(http.StatusOK, players)
}

// GetServerLogs はサーバーの直近のログを配列形式で返すエンドポイント（GET /api/server/logs）です。
func (h *MinecraftHandler) GetServerLogs(c echo.Context) error {
	lines, err := h.fileService.GetServerLogs()
	if err != nil {
		return c.JSON(http.StatusOK, []string{}) // エラー時は空で返す
	}
	return c.JSON(http.StatusOK, lines)
}
