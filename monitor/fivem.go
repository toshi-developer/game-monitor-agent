package monitor

import (
	"fmt"
	"net/http"
	"time"
)

func CheckFiveM(name, addr string, timeout time.Duration, start time.Time) Result {
	url := fmt.Sprintf("http://%s/info.json", addr)
	client := http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("[DEBUG] [%s] FiveM API応答なし: %v\n", name, err)
		return Result{Name: name, IsAlive: false, Message: "API Unreachable"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Name: name, IsAlive: false, Message: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	fmt.Printf("[DEBUG] [%s] FiveM API応答 OK\n", name)
	return Result{Name: name, IsAlive: true, Latency: time.Since(start), Message: "FiveM Fully Functional"}
}
