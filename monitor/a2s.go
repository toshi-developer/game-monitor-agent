package monitor

import (
	"fmt"
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
		fmt.Printf("[DEBUG] [A2S] %s への A2S_INFO クエリ失敗: %v\n", addr, err)
		return res
	}

	res.IsAlive = true
	res.Latency = time.Since(start)
	res.PlayerCount = int(info.Players)
	res.MaxPlayers = int(info.MaxPlayers)
	res.MapName = info.Map
	res.Version = info.Version
	res.Message = fmt.Sprintf("A2S Active: %s", info.Name)

	fmt.Printf("[DEBUG] [A2S] %s 監視完了: Players=%d/%d, Map=%s, Version=%s, Latency=%v\n",
		addr, res.PlayerCount, res.MaxPlayers, res.MapName, res.Version, res.Latency)

	return res
}
