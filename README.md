# Game Monitor Agent (by toshi dev)

Go言語で開発された、軽量で強力なゲームサーバ監視エージェントです。
特に **FiveM** サーバの死活監視に最適化されており、詳細なメトリクスを時系列DB（InfluxDB）に蓄積します。



## 主な機能
- **並列監視**: 複数のゲームサーバを Goroutine で同時にチェック。
- **FiveM 特化型ロジック**: ポート疎通(L4)だけでなく、`/info.json` へのHTTP応答(L7)まで確認。
- **時系列DB連携**: InfluxDB 2.x への自動データ送信。
- **デバッグフレンドリー**: `[DEBUG]` プレフィックス付きの詳細なログ出力。
- **常駐化対応**: Docker および systemd による 24/365 の動作をサポート。

## クイックスタート

### 1. インフラの起動
Dockerを使用して InfluxDB と Grafana を起動します。
```bash
docker-compose up -d
```

### 2. 設定
`config.yaml.example` を `config.yaml` にコピーし、環境に合わせて編集してください。
```bash
cp config.yaml.example config.yaml
```

### 3. エージェントの実行
```bash
go mod tidy
go run main.go
```

## デプロイと常駐化 (Deployment & Persistence)

### 1. Docker Compose による常駐化 (推奨)
`docker-compose.yml` 内の `restart: always` 設定により、OS起動時やエラー終了時にエージェントが自動で再起動します。

**起動方法:**
```bash
docker-compose up -d --build
```

**ログの確認方法:**
```bash
docker logs -f game-monitor-agent
```

---

### 2. systemd による常駐化 (Linux専用)
Dockerを使用せず、Linuxのシステムサービスとして直接実行する場合の設定例です。

1. バイナリをビルドします:
   ```bash
   go build -o agent main.go
   ```
2. `/etc/systemd/system/game-monitor.service` を作成します（パスやユーザ名は環境に合わせて変更してください）:
   ```ini
   [Unit]
   Description=Game Monitor Agent (toshi dev)
   After=network.target influxdb.service

   [Service]
   Type=simple
   User=YOUR_USER_NAME
   WorkingDirectory=/path/to/game-monitor-agent
   ExecStart=/path/to/game-monitor-agent/agent
   Restart=always
   RestartSec=10

   [Install]
   WantedBy=multi-user.target
   ```
3. サービスを有効化して起動します:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable game-monitor.service
   sudo systemctl start game-monitor.service
   ```

## 設定項目 (config.yaml)
| 項目 | 説明 |
| :--- | :--- |
| `game_type` | `fivem` または `generic` を指定可能。 |
| `interval_seconds` | 監視を実行する間隔（秒）。 |
| `timeout_ms` | 応答待ちのタイムアウト（ミリ秒）。 |

## 開発環境
- **Language**: Go 1.21+
- **Database**: InfluxDB 2.x
- **Visualization**: Grafana

## ライセンス
MIT License

## 開発者
- ユーザ名: **toshi-developer**
- 屋号: **toshi dev**
- Web: [GitHub Profile](https://github.com/toshi-developer)