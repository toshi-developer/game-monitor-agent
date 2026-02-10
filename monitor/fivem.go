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
	resp, err := client.Get(fmt.Sprintf("http://%s/dynamic.json", addr))
	if err == nil {
		defer resp.Body.Close()
		var dyn FiveMDynamic
		if err := json.NewDecoder(resp.Body).Decode(&dyn); err == nil {
			res.PlayerCount = dyn.Clients
			res.MaxPlayers = dyn.MaxClients
		}
	}

	// 2. 死活確認とレイテンシ
	respInfo, err := client.Get(fmt.Sprintf("http://%s/info.json", addr))
	if err != nil {
		res.IsAlive = false
		res.Message = "API Unreachable"
		return res
	}
	defer respInfo.Body.Close()

	res.IsAlive = (respInfo.StatusCode == http.StatusOK)
	res.Latency = time.Since(start)
	res.Message = "FiveM Metrics Collected"

	return res
}
