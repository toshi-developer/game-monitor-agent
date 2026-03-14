package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/toshi-developer/game-monitor-agent/config"
)

// providers はゲーム種別名から Provider へのレジストリです。
// init() 関数で各プロバイダが自身を登録します。
var providers = map[string]Provider{}

// RegisterProvider はプロバイダをレジストリに登録します。
func RegisterProvider(gameType string, p Provider) {
	providers[gameType] = p
}

// RunAll は全サーバーの監視を並列実行し、結果を返します。
// ホストのシステムメトリクスは1サイクルにつき1回だけ取得します。
func RunAll(servers []config.ServerConfig) []Result {
	// 1. システムメトリクスを1回だけ取得
	sysMetrics := collectSystemMetrics()

	// 2. 各サーバーのゲーム固有監視を並列実行
	var wg sync.WaitGroup
	results := make(chan Result, len(servers))

	for _, s := range servers {
		wg.Add(1)
		go func(c config.ServerConfig) {
			defer wg.Done()
			results <- execute(c, sysMetrics)
		}(s)
	}

	wg.Wait()
	close(results)

	var output []Result
	for r := range results {
		output = append(output, r)
	}
	return output
}

// execute は単一サーバーの監視を実行します。
func execute(c config.ServerConfig, sysMetrics SystemMetrics) Result {
	timeout := time.Duration(c.TimeoutMS) * time.Millisecond
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)

	res := Result{
		Name:          c.Name,
		SystemMetrics: sysMetrics,
	}

	p, ok := providers[c.GameType]
	if !ok {
		// 未登録のゲーム種別
		fmt.Printf("[WARN] [%s] 未知のゲーム種別: %q (スキップ)\n", c.Name, c.GameType)
		res.GameResult = GameResult{
			IsAlive: false,
			Message: fmt.Sprintf("Unknown game_type: %q", c.GameType),
		}
		return res
	}

	// 7DtD はWeb APIアクセスのために ServerConfig を渡す
	if sp, ok := p.(*SevenDTDProvider); ok {
		res.GameResult = sp.FetchWithWebAPI(addr, timeout, &c)
	} else {
		res.GameResult = p.Fetch(addr, timeout)
	}
	return res
}

// collectSystemMetrics はホストのシステムリソース情報を1回取得します。
func collectSystemMetrics() SystemMetrics {
	var m SystemMetrics

	// CPU
	c, err := cpu.Percent(0, false)
	if err == nil && len(c) > 0 {
		m.CPUUsage = c[0]
	} else {
		fmt.Printf("[DEBUG] CPU使用率の取得に失敗しました: %v\n", err)
	}

	// Memory
	vm, err := mem.VirtualMemory()
	if err == nil {
		m.MemUsage = vm.UsedPercent
	}
	sm, err := mem.SwapMemory()
	if err == nil {
		m.SwapUsage = sm.UsedPercent
	}

	// Disk
	d, err := disk.Usage("/")
	if err == nil {
		m.DiskUsage = d.UsedPercent
	} else {
		fmt.Printf("[DEBUG] ディスク使用率の取得に失敗しました: %v\n", err)
	}

	// Network
	io, err := net.IOCounters(false)
	if err == nil && len(io) > 0 {
		m.NetSent = io[0].BytesSent / 1024
		m.NetRecv = io[0].BytesRecv / 1024
	}

	// Connections
	conns, err := net.Connections("tcp")
	if err == nil {
		count := 0
		for _, conn := range conns {
			if conn.Status == "ESTABLISHED" {
				count++
			}
		}
		m.Connections = count
	}

	fmt.Printf("[DEBUG] リソース取得完了: CPU:%.1f%%, Mem:%.1f%%, Disk:%.1f%%, NetSent:%dKB, NetRecv:%dKB, Conns:%d\n",
		m.CPUUsage, m.MemUsage, m.DiskUsage, m.NetSent, m.NetRecv, m.Connections)

	return m
}
