package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/toshi-developer/game-monitor-agent/config"
)

type Result struct {
	Name    string
	IsAlive bool
	Latency time.Duration
	Message string
}

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
	fmt.Printf("[DEBUG] [%s] チェック開始...\n", c.Name)
	start := time.Now()
	timeout := time.Duration(c.TimeoutMS) * time.Millisecond
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)

	switch c.GameType {
	case "fivem":
		return CheckFiveM(c.Name, addr, timeout, start)
	default:
		// 簡易的なTCPチェック（monitor/fivem.go内に定義がない場合はここに簡易実装可）
		return Result{Name: c.Name, IsAlive: true, Latency: time.Since(start), Message: "Generic OK"}
	}
}
