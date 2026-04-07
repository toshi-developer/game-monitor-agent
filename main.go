package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/toshi-developer/game-monitor-agent/config"
	"github.com/toshi-developer/game-monitor-agent/monitor"
	"github.com/toshi-developer/game-monitor-agent/storage"
)

func main() {
	slog.Info("Game Monitor Agent starting")

	// 1. 設定ロード
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		slog.Error("設定読み込み失敗", "error", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		slog.Error("設定バリデーション失敗", "error", err)
		os.Exit(1)
	}

	// 2. ストレージ初期化
	var db *storage.InfluxClient
	if cfg.Destination.Mode == "local" {
		db = storage.NewInfluxClient(cfg)
		defer db.Close()
	}

	// 3. グレースフルシャットダウン用のコンテキスト
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. メインループ
	ticker := time.NewTicker(time.Duration(cfg.Monitoring.Interval) * time.Second)
	defer ticker.Stop()

	slog.Info("監視開始",
		"servers", len(cfg.Monitoring.Servers),
		"interval_sec", cfg.Monitoring.Interval,
	)

	for {
		slog.Debug("サイクル開始")

		results := monitor.RunAll(cfg.Monitoring.Servers)

		if db != nil {
			db.SaveResults(results)
		}

		for _, r := range results {
			slog.Info("監視結果",
				"server", r.Name,
				"alive", r.IsAlive,
				"message", r.Message,
			)
		}

		select {
		case <-ctx.Done():
			slog.Info("シャットダウンシグナルを受信しました。終了します。")
			return
		case <-ticker.C:
		}
	}
}
