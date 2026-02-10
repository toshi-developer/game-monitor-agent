package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FiveMDynamic struct {
	Clients    int `json:"clients"`
	MaxClients int `json:"sv_maxclients"`
}

func fetchFiveMDetails(res Result, addr string, timeout time.Duration, start time.Time) Result {
	client := http.Client{Timeout: timeout}

	// 1. プレイヤー人数の取得
	dynamicURL := fmt.Sprintf("http://%s/dynamic.json", addr)
	fmt.Printf("[DEBUG] [%s] FiveM Dynamic APIへリクエスト: %s\n", res.Name, dynamicURL)
	resp, err := client.Get(dynamicURL)
	if err == nil {
		defer resp.Body.Close()
		var dyn FiveMDynamic
		if err := json.NewDecoder(resp.Body).Decode(&dyn); err == nil {
			res.PlayerCount = dyn.Clients
			res.MaxPlayers = dyn.MaxClients
			fmt.Printf("[DEBUG] [%s] プレイヤー数取得成功: %d/%d\n", res.Name, res.PlayerCount, res.MaxPlayers)
		} else {
			fmt.Printf("[DEBUG] [%s] dynamic.json のパースに失敗しました\n", res.Name)
		}
	} else {
		fmt.Printf("[DEBUG] [%s] dynamic.json の取得に失敗しました: %v\n", res.Name, err)
	}

	// 2. 死活確認とレイテンシ
	infoURL := fmt.Sprintf("http://%s/info.json", addr)
	fmt.Printf("[DEBUG] [%s] FiveM Info APIへリクエスト: %s\n", res.Name, infoURL)
	respInfo, err := client.Get(infoURL)
	if err != nil {
		res.IsAlive = false
		res.Message = fmt.Sprintf("API Unreachable: %v", err)
		fmt.Printf("[DEBUG] [%s] Info API接続エラー: %v\n", res.Name, err)
		return res
	}
	defer respInfo.Body.Close()

	res.IsAlive = (respInfo.StatusCode == http.StatusOK)
	res.Latency = time.Since(start)
	res.Message = fmt.Sprintf("Status: %d, Latency: %v", respInfo.StatusCode, res.Latency)

	fmt.Printf("[DEBUG] [%s] FiveM監視完了: Alive:%v, Latency:%v\n", res.Name, res.IsAlive, res.Latency)

	return res
}
