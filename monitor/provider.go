package monitor

import "time"

// Provider はゲーム固有の監視ロジックを実装するインターフェースです。
// 各ゲームプロバイダはこのインターフェースを満たすことで、エンジンに登録できます。
type Provider interface {
	// Fetch はゲームサーバーに問い合わせを行い、GameResult を返します。
	// addr は "host:port" 形式のアドレス、timeout は接続タイムアウトです。
	Fetch(addr string, timeout time.Duration) GameResult
}

// GameResult はゲーム固有の監視結果です（プロバイダが返す値）。
type GameResult struct {
	IsAlive     bool
	Latency     time.Duration
	Message     string
	PlayerCount int
	MaxPlayers  int
	Version     string
	MapName     string
	GameTime    string // 7DtD 等のゲーム内時間 ("Day 7 21:00" 形式)
}

// SystemMetrics はホストシステムのリソース情報です（全サーバー共通、1サイクルに1回取得）。
type SystemMetrics struct {
	CPUUsage    float64
	MemUsage    float64
	SwapUsage   float64
	DiskUsage   float64
	NetSent     uint64 // KB
	NetRecv     uint64 // KB
	Connections int
}

// Result は監視結果の統一フォーマットです。GameResult + SystemMetrics を合成したものです。
type Result struct {
	Name string
	GameResult
	SystemMetrics
}
