package storage

import (
	"context"
	"log/slog"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/toshi-developer/game-monitor-agent/config"
	"github.com/toshi-developer/game-monitor-agent/monitor"
)

type InfluxClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	bucket   string
	org      string
}

func NewInfluxClient(cfg *config.Config) *InfluxClient {
	c := cfg.Destination.Local
	client := influxdb2.NewClient(c.URL, c.Token)
	writeAPI := client.WriteAPIBlocking(c.Org, c.Bucket)
	return &InfluxClient{client: client, writeAPI: writeAPI, bucket: c.Bucket, org: c.Org}
}

func (ic *InfluxClient) SaveResults(results []monitor.Result) {
	for _, res := range results {
		status := 0
		if res.IsAlive {
			status = 1
		}

		log := slog.With("server", res.Name)

		fields := map[string]interface{}{
			"is_alive":    status,
			"latency_ms":  res.Latency.Milliseconds(),
			"cpu_usage":   res.CPUUsage,
			"mem_usage":   res.MemUsage,
			"swap_usage":  res.SwapUsage,
			"disk_usage":  res.DiskUsage,
			"net_sent_kb": res.NetSent,
			"net_recv_kb": res.NetRecv,
			"connections": res.Connections,
			"players":     res.PlayerCount,
			"max_players": res.MaxPlayers,
		}

		// ゲーム固有フィールド（値がある場合のみ書き込み）
		if res.MapName != "" {
			fields["map_name"] = res.MapName
		}
		if res.Version != "" {
			fields["version"] = res.Version
		}
		if res.GameTime != "" {
			fields["game_time"] = res.GameTime
		}

		p := influxdb2.NewPoint("server_metrics",
			map[string]string{"server_name": res.Name},
			fields,
			time.Now())

		log.Debug("InfluxDB へデータ送信中", "bucket", ic.bucket)

		if err := ic.writeAPI.WritePoint(context.Background(), p); err != nil {
			log.Error("InfluxDB 書き込み失敗", "error", err)
		} else {
			log.Debug("InfluxDB 書き込み成功")
		}
	}
}

func (ic *InfluxClient) Close() { ic.client.Close() }
