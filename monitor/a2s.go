package monitor

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/rumblefrog/go-a2s"
)

// A2SProvider は Steam Source Query (A2S) を使用するゲームサーバーの汎用プロバイダです。
// ARK: Survival Ascended など、A2S プロトコルに対応するゲームで使用します。
type A2SProvider struct{}

func init() {
	RegisterProvider("ark", &A2SProvider{})
}

func (p *A2SProvider) Fetch(addr string, timeout time.Duration) GameResult {
	log := slog.With("provider", "a2s", "addr", addr)
	var res GameResult

	start := time.Now()
	client, err := a2s.NewClient(addr, a2s.TimeoutOption(timeout))
	if err != nil {
		res.IsAlive = false
		res.Message = fmt.Sprintf("A2S Client Error: %v", err)
		return res
	}
	defer client.Close()

	info, err := client.QueryInfo()
	if err != nil {
		res.IsAlive = false
		res.Message = fmt.Sprintf("A2S Query Failed: %v", err)
		log.Warn("A2S_INFO クエリ失敗", "error", err)
		return res
	}

	res.IsAlive = true
	res.Latency = time.Since(start)
	res.PlayerCount = int(info.Players)
	res.MaxPlayers = int(info.MaxPlayers)
	res.MapName = info.Map
	res.Version = info.Version
	res.Message = fmt.Sprintf("A2S Active: %s", info.Name)

	log.Debug("監視完了",
		"players", res.PlayerCount,
		"max_players", res.MaxPlayers,
		"map", res.MapName,
		"version", res.Version,
		"latency", res.Latency,
	)

	return res
}
