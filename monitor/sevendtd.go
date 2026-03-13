package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rumblefrog/go-a2s"
	"github.com/toshi-developer/game-monitor-agent/config"
)

// SevenDTDProvider は 7 Days to Die サーバーの監視プロバイダです。
// A2S による基本情報取得に加え、Web API からゲーム内時間を取得します。
type SevenDTDProvider struct{}

func init() {
	RegisterProvider("7dtd", &SevenDTDProvider{})
}

// SevenDTDStats は 7DtD Web API (/api/getstats) のレスポンス構造です。
type SevenDTDStats struct {
	GameTime struct {
		Days    int `json:"days"`
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
	} `json:"gameTime"`
}

func (p *SevenDTDProvider) Fetch(addr string, timeout time.Duration) GameResult {
	return p.FetchWithWebAPI(addr, timeout, nil)
}

// FetchWithWebAPI は ServerConfig を受け取り、Web API ポートを使ってゲーム内時間も取得します。
func (p *SevenDTDProvider) FetchWithWebAPI(addr string, timeout time.Duration, serverCfg *config.ServerConfig) GameResult {
	var res GameResult

	// 1. A2S クエリによる基本情報取得
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
		fmt.Printf("[DEBUG] [7DtD] %s への A2S_INFO クエリ失敗: %v\n", addr, err)
		return res
	}

	res.IsAlive = true
	res.Latency = time.Since(start)
	res.PlayerCount = int(info.Players)
	res.MaxPlayers = int(info.MaxPlayers)
	res.MapName = info.Map
	res.Version = info.Version
	res.Message = fmt.Sprintf("7DtD Active: %s", info.Name)

	// 2. Web API からゲーム内時間を取得
	if serverCfg != nil && serverCfg.WebAPIPort > 0 {
		gameTime, err := fetchGameTime(serverCfg.Address, serverCfg.WebAPIPort, timeout)
		if err != nil {
			fmt.Printf("[DEBUG] [7DtD] ゲーム内時間の取得に失敗: %v\n", err)
		} else {
			res.GameTime = gameTime
		}
	}

	fmt.Printf("[DEBUG] [7DtD] %s 監視完了: Players=%d/%d, Map=%s, GameTime=%s, Latency=%v\n",
		addr, res.PlayerCount, res.MaxPlayers, res.MapName, res.GameTime, res.Latency)

	return res
}

// fetchGameTime は 7DtD の Web API からゲーム内時間を取得し、"Day 7 21:00" 形式で返します。
func fetchGameTime(host string, webAPIPort int, timeout time.Duration) (string, error) {
	url := fmt.Sprintf("http://%s:%d/api/getstats", host, webAPIPort)
	client := http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("Web API 通信エラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Web API ステータス異常: %d", resp.StatusCode)
	}

	var stats SevenDTDStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return "", fmt.Errorf("Web API パースエラー: %w", err)
	}

	return fmt.Sprintf("Day %d %02d:%02d", stats.GameTime.Days, stats.GameTime.Hours, stats.GameTime.Minutes), nil
}
