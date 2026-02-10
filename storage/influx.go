package storage

import (
	"context"
	"fmt"
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

		p := influxdb2.NewPoint("server_metrics",
			map[string]string{"server_name": res.Name},
			map[string]interface{}{
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
			},
			time.Now())

		if err := ic.writeAPI.WritePoint(context.Background(), p); err != nil {
			fmt.Printf("[DEBUG] [%s] DB Error: %v\n", res.Name, err)
		}
	}
}

func (ic *InfluxClient) Close() { ic.client.Close() }
