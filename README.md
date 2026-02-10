## デプロイと常駐化 (Deployment & Persistence)

本エージェントをサーバ上で24時間365日安定して動作させるための推奨設定です。

### 1. Docker Compose による常駐化 (推奨)
`docker-compose.yml` 内の `restart: always` 設定により、OS起動時やエラー終了時にエージェントが自動で再起動します。

**実行方法:**
```bash
# バックグラウンドで起動
docker-compose up -d
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
2. `/etc/systemd/system/game-monitor.service` を以下の内容で作成します:
   ```ini
   [Unit]
   Description=Game Monitor Agent (toshi dev)
   After=network.target

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

---

## メンテナンスとデバッグ

すべてのパッケージ（`config`, `monitor`, `storage`）において詳細なデバッグメッセージを出力するように設計されています。

`[DEBUG]` プレフィックスがついたログを確認することで、設定の読み込み状況、FiveM APIへの疎通、InfluxDBへの書き込み成功可否をリアルタイムで追跡可能です。