package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FiveMDynamic の sv_maxclients に ,string を追加して
// 文字列で返ってくる数値に対応させます
type FiveMDynamic struct {
	Clients    int `json:"clients"`
	MaxClients int `json:"sv_maxclients,string"`
}

func fetchFiveMDetails(res Result, addr string, timeout time.Duration, start time.Time) Result {
	client := http.Client{Timeout: timeout}

	// 1. プレイヤー人数と最大人数の取得 (dynamic.json)
	dynURL := fmt.Sprintf("http://%s/dynamic.json", addr)
	fmt.Printf("[DEBUG] [%s] FiveM Dynamic APIへリクエスト開始: %s\n", res.Name, dynURL)

	resp, err := client.Get(dynURL)
	if err == nil {
		defer resp.Body.Close()
		var dyn FiveMDynamic
		if err := json.NewDecoder(resp.Body).Decode(&dyn); err == nil {
			res.PlayerCount = dyn.Clients
			res.MaxPlayers = dyn.MaxClients
			fmt.Printf("[DEBUG] [%s] プレイヤー数取得成功: %d/%d (JSONパース完了)\n", res.Name, res.PlayerCount, res.MaxPlayers)
		} else {
			// 具体的なパースエラー内容を出力
			fmt.Printf("[DEBUG] [%s] dynamic.json のパースに失敗しました: %v\n", res.Name, err)
			fmt.Printf("[DEBUG] [%s] ヒント: sv_maxclients が文字列形式であることを確認済みです\n", res.Name)
		}
	} else {
		fmt.Printf("[DEBUG] [%s] dynamic.json の通信エラー: %v\n", res.Name, err)
	}

	// 2. 死活確認とレイテンシの取得 (info.json)
	infoURL := fmt.Sprintf("http://%s/info.json", addr)
	fmt.Printf("[DEBUG] [%s] FiveM Info APIへリクエスト開始: %s\n", res.Name, infoURL)

	respInfo, err := client.Get(infoURL)
	if err != nil {
		res.IsAlive = false
		res.Message = fmt.Sprintf("API Unreachable: %v", err)
		fmt.Printf("[DEBUG] [%s] Info API接続失敗: %v\n", res.Name, err)
		return res
	}
	defer respInfo.Body.Close()

	res.IsAlive = (respInfo.StatusCode == http.StatusOK)
	res.Latency = time.Since(start)
	res.Message = fmt.Sprintf("FiveM Active (Status: %d)", respInfo.StatusCode)

	fmt.Printf("[DEBUG] [%s] FiveM監視全プロセス完了: Alive=%v, Latency=%v\n", res.Name, res.IsAlive, res.Latency)

	return res
}
