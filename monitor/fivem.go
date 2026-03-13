package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FiveMProvider は FiveM サーバーの監視プロバイダです。
type FiveMProvider struct{}

func init() {
	RegisterProvider("fivem", &FiveMProvider{})
}

// FiveMDynamic の sv_maxclients に ,string を追加して
// 文字列で返ってくる数値に対応させます
type FiveMDynamic struct {
	Clients    int `json:"clients"`
	MaxClients int `json:"sv_maxclients,string"`
}

func (p *FiveMProvider) Fetch(addr string, timeout time.Duration) GameResult {
	client := http.Client{Timeout: timeout}
	var res GameResult

	// 1. プレイヤー人数と最大人数の取得 (dynamic.json)
	if err := p.fetchDynamic(&client, addr, &res); err != nil {
		fmt.Printf("[DEBUG] [FiveM] dynamic.json の取得に失敗: %v\n", err)
	}

	// 2. 死活確認とレイテンシの取得 (info.json)
	start := time.Now()
	infoURL := fmt.Sprintf("http://%s/info.json", addr)
	resp, err := client.Get(infoURL)
	if err != nil {
		res.IsAlive = false
		res.Message = fmt.Sprintf("API Unreachable: %v", err)
		return res
	}
	resp.Body.Close()

	res.IsAlive = (resp.StatusCode == http.StatusOK)
	res.Latency = time.Since(start)
	res.Message = fmt.Sprintf("FiveM Active (Status: %d)", resp.StatusCode)

	fmt.Printf("[DEBUG] [FiveM] 監視完了: Alive=%v, Players=%d/%d, Latency=%v\n",
		res.IsAlive, res.PlayerCount, res.MaxPlayers, res.Latency)

	return res
}

func (p *FiveMProvider) fetchDynamic(client *http.Client, addr string, res *GameResult) error {
	dynURL := fmt.Sprintf("http://%s/dynamic.json", addr)
	resp, err := client.Get(dynURL)
	if err != nil {
		return fmt.Errorf("通信エラー: %w", err)
	}
	defer resp.Body.Close()

	var dyn FiveMDynamic
	if err := json.NewDecoder(resp.Body).Decode(&dyn); err != nil {
		return fmt.Errorf("パースエラー: %w", err)
	}

	res.PlayerCount = dyn.Clients
	res.MaxPlayers = dyn.MaxClients
	fmt.Printf("[DEBUG] [FiveM] プレイヤー数取得成功: %d/%d\n", res.PlayerCount, res.MaxPlayers)
	return nil
}
