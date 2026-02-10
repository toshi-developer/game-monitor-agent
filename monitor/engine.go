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

func RunAll(servers []config.ServerConfig) []Result {
	var wg sync.WaitGroup
	results := make(chan Result, len(servers))

	for _, s := range servers {
		wg.Add(1)
		go func(c config.ServerConfig) {
			defer wg.Done()
			results <- execute(c)
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

func execute(c config.ServerConfig) Result {
	start := time.Now()
	timeout := time.Duration(c.TimeoutMS) * time.Millisecond
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)

	// 基盤となるリソース情報を取得
	res := Result{Name: c.Name}
	fillSystemMetrics(&res)

	// ゲーム種別ごとの詳細取得
	switch c.GameType {
	case "fivem":
		return fetchFiveMDetails(res, addr, timeout, start)
	default:
		// 未知のゲームはTCP疎通確認のみ（仮）
		res.IsAlive = true // 実際にはここで簡単な疎通確認を入れる
		res.Message = "Generic Check"
		return res
	}
}

func fillSystemMetrics(res *Result) {
	// CPU
	c, err := cpu.Percent(0, false)
	if err == nil && len(c) > 0 {
		res.CPUUsage = c[0]
	} else {
		fmt.Printf("[DEBUG] [%s] CPU使用率の取得に失敗しました: %v\n", res.Name, err)
	}

	// Memory
	vm, err := mem.VirtualMemory()
	if err == nil {
		res.MemUsage = vm.UsedPercent
	}
	sm, err := mem.SwapMemory()
	if err == nil {
		res.SwapUsage = sm.UsedPercent
	}

	// Disk
	d, err := disk.Usage("/")
	if err == nil {
		res.DiskUsage = d.UsedPercent
	} else {
		fmt.Printf("[DEBUG] [%s] ディスク使用率の取得に失敗しました: %v\n", res.Name, err)
	}

	// Network
	io, err := net.IOCounters(false)
	if err == nil && len(io) > 0 {
		res.NetSent = io[0].BytesSent / 1024
		res.NetRecv = io[0].BytesRecv / 1024
	}

	// Connections
	conns, err := net.Connections("tcp")
	count := 0
	if err == nil {
		for _, conn := range conns {
			if conn.Status == "ESTABLISHED" {
				count++
			}
		}
		res.Connections = count
	}

	// 詳細な取得結果をログ出力
	fmt.Printf("[DEBUG] [%s] リソース取得完了: CPU:%.1f%%, Mem:%.1f%%, Disk:%.1f%%, NetSent:%dKB, NetRecv:%dKB, Conns:%d\n",
		res.Name, res.CPUUsage, res.MemUsage, res.DiskUsage, res.NetSent, res.NetRecv, res.Connections)
}
