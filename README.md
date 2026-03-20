# Minecraft Server Manager Web App

ローカル環境での稼働を前提とした、マインクラフトサーバー管理用のWEBアプリケーションです。

## 主な機能
- **サーバーの起動、停止、再起動**
- **サーバーの参加メンバー一覧の表示**
- **参加ユーザーの管理** (キック、バン、OP付与・剥奪)
- **ホワイトリスト管理**
- **ホワイトモード切り替え**
- **サーバーログの確認**
- **コマンド実行**

## 技術スタック
- **Frontend**: SvelteKit + TypeScript + TailwindCSS
- **Backend**: Go + Echo
- **Database**: SQLite

## デザインテーマ
- **Vaporwave** (ネオンピンク、シアン、パープルを基調としたレトロフューチャーなデザイン)

## プロジェクト構成 (推奨)
```
.
├── frontend/       # SvelteKitフロントエンドアプリケーション
├── backend/        # Goバックエンドアプリケーション
├── server/         # ユーザーが用意したマインクラフトサーバーのjarファイルを配置するディレクトリ
│   └── server.jar  # 汎用的なマインクラフトサーバーの実行ファイル
└── data/           # SQLiteのデータベースファイル(e.g., mcmt.db)の保存場所
```

## セットアップと起動方法

### 前提条件
- Node.js (SvelteKit用)
- Go言語環境
- SQLite

**注意**: マインクラフトサーバーの動作には別途Javaが必要ですが、本アプリのセットアップ手順には含まれません。ユーザー自身で適切なバージョンのJavaと、実行したいマインクラフトサーバーの`jar`ファイル（Vanilla, Paper, Spigot等汎用的に対応）をご用意ください。

### バックエンド (Go)
1. `backend` ディレクトリに移動します。
2. 依存関係をインストールします。
   ```bash
   go mod download
   ```
3. アプリケーションをビルド・起動します。
   ```bash
   go build -o server-manager
   ./server-manager
   ```

### フロントエンド (SvelteKit)
1. `frontend` ディレクトリに移動します。
2. 依存関係をインストールします。
   ```bash
   npm install
   ```
3. 開発サーバーを起動するか、ビルドして本番用に起動します。
   ```bash
   # 開発用
   npm run dev

   # ビルド用
   npm run build
   node build
   ```

## サーバーファイルの設定
ユーザーが用意したマインクラフトサーバーの `jar` ファイルは `server/server.jar` として配置してください。アプリはデフォルトでこのパスを参照してサーバーを起動します。
