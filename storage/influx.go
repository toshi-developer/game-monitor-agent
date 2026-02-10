package storage

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api" // ← これを追加
	"github.com/toshi-developer/game-monitor-agent/config"
	"github.com/toshi-developer/game-monitor-agent/monitor"
)

type InfluxClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking // ← ここを api.WriteAPIBlocking に修正
	bucket   string
	org      string
}

func NewInfluxClient(cfg *config.Config) *InfluxClient {
	c := cfg.Destination.Local
	fmt.Printf("[DEBUG] InfluxDBに接続中: %s\n", c.URL)

	client := influxdb2.NewClient(c.URL, c.Token)
	// クライアントから WriteAPIBlocking を取得
	writeAPI := client.WriteAPIBlocking(c.Org, c.Bucket)

	return &InfluxClient{
		client:   client,
		writeAPI: writeAPI,
		bucket:   c.Bucket,
		org:      c.Org,
	}
}

func (ic *InfluxClient) SaveResults(results []monitor.Result) {
	for _, res := range results {
		status := 0
		if res.IsAlive {
			status = 1
		}

		// ポイントの作成
		p := influxdb2.NewPoint("server_metrics",
			map[string]string{"server_name": res.Name},
			map[string]interface{}{
				"is_alive": status,
				"latency":  res.Latency.Milliseconds(),
			},
			time.Now())

		// 書き込み実行
		if err := ic.writeAPI.WritePoint(context.Background(), p); err != nil {
			fmt.Printf("[DEBUG] [%s] DB保存失敗: %v\n", res.Name, err)
		} else {
			fmt.Printf("[DEBUG] [%s] メトリクスをInfluxDBに保存しました\n", res.Name)
		}
	}
}

func (ic *InfluxClient) Close() {
	ic.client.Close()
}
