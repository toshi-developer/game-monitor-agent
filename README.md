# Game Monitor Agent (by toshi dev)

Go言語で開発された、軽量で強力なマルチゲームサーバ監視エージェントです。
FiveMをはじめとするゲームサーバの死活監視に加え、ホストマシンのリソース状況をリアルタイムで可視化します。



## 主な機能
- **マルチゲーム対応**: インターフェース設計により、FiveM以外のゲームも容易に拡張可能。
- **フルスタック監視**: 
  - **Game**: ステータス、プレイヤー人数、最大人数、レイテンシ。
  - **System**: CPU使用率、メモリ使用率、ディスク使用率、ネットワークI/O、コネクション数。
- **時系列DB連携**: InfluxDB 2.x への自動データ送信。
- **デバッグ機能**: `[DEBUG]` ログにより、システム情報とゲーム情報の取得プロセスを可視化。

## クイックスタート

### 1. インフラの起動
```bash
docker-compose up -d
```

### 2. エージェントのビルドと起動
```bash
# 依存関係の整理
go mod tidy
# コンテナの再ビルドと起動
docker-compose up -d --build
```

## 監視メトリクス一覧 (InfluxDB Fields)
| フィールド | 内容 |
| :--- | :--- |
| `is_alive` | サーバーの生存状態 (1: Online, 0: Offline) |
| `players` | 現在の接続プレイヤー数 |
| `cpu_usage` | ホストマシンのCPU使用率 (%) |
| `mem_usage` | ホストマシンのメモリ使用率 (%) |
| `net_recv_kb` | ネットワーク受信量 (KB) |

## 開発者
- **toshi-developer** (toshi dev)
- [GitHub Profile](https://github.com/toshi-developer)