package monitor

import "time"

// Result は監視結果の統一フォーマットです
type Result struct {
	Name    string
	IsAlive bool
	Latency time.Duration
	Message string
	// システムリソース (全ゲーム共通)
	CPUUsage    float64
	MemUsage    float64
	SwapUsage   float64
	DiskUsage   float64 // ← ここが不足していた可能性があります
	NetSent     uint64  // KB
	NetRecv     uint64  // KB
	Connections int
	// ゲーム内情報 (対応しているゲームのみ)
	PlayerCount int
	MaxPlayers  int
	Version     string
}
