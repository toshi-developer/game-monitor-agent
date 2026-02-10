package main

import (
	"fmt"
	"log"
	"time"

	"github.com/toshi-developer/game-monitor-agent/config"
	"github.com/toshi-developer/game-monitor-agent/monitor"
	"github.com/toshi-developer/game-monitor-agent/storage"
)

func main() {
	fmt.Println("=== Toshi Dev: Game Monitor Agent Starting ===")

	// 1. 設定ロード
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("[FATAL] 設定読み込み失敗: %v", err)
	}

	// 2. ストレージ初期化
	var db *storage.InfluxClient
	if cfg.Destination.Mode == "local" {
		db = storage.NewInfluxClient(cfg)
		defer db.Close()
	}

	// 3. メインループ
	ticker := time.NewTicker(time.Duration(cfg.Monitoring.Interval) * time.Second)
	defer ticker.Stop()

	for {
		fmt.Printf("\n[DEBUG] サイクル開始: %s\n", time.Now().Format("15:04:05"))

		results := monitor.RunAll(cfg.Monitoring.Servers)

		if db != nil {
			db.SaveResults(results)
		}

		for _, r := range results {
			fmt.Printf("[RESULT] %-15s | Alive: %-5v | Msg: %s\n", r.Name, r.IsAlive, r.Message)
		}

		<-ticker.C
	}
}
