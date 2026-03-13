package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Name       string `yaml:"name"`
	GameType   string `yaml:"game_type"`
	Address    string `yaml:"address"`
	Port       int    `yaml:"port"`
	TimeoutMS  int    `yaml:"timeout_ms"`
	WebAPIPort int    `yaml:"web_api_port,omitempty"` // 7DtD Web API 用ポート (デフォルト: 0=無効)
}

type Config struct {
	Monitoring struct {
		Interval int            `yaml:"interval_seconds"`
		Servers  []ServerConfig `yaml:"game_servers"`
	} `yaml:"monitoring"`
	Destination struct {
		Mode  string `yaml:"mode"`
		Local struct {
			URL    string `yaml:"url"`
			Token  string `yaml:"token"`
			Org    string `yaml:"org"`
			Bucket string `yaml:"bucket"`
		} `yaml:"local"`
	} `yaml:"destination"`
}

// Validate は設定値のバリデーションを行います。
func (c *Config) Validate() error {
	if c.Monitoring.Interval <= 0 {
		return fmt.Errorf("interval_seconds は1以上である必要があります (現在: %d)", c.Monitoring.Interval)
	}
	if len(c.Monitoring.Servers) == 0 {
		return fmt.Errorf("game_servers が空です。最低1つのサーバーを設定してください")
	}
	for i, s := range c.Monitoring.Servers {
		if s.Name == "" {
			return fmt.Errorf("game_servers[%d]: name が空です", i)
		}
		if s.Address == "" {
			return fmt.Errorf("game_servers[%d] (%s): address が空です", i, s.Name)
		}
		if s.Port <= 0 || s.Port > 65535 {
			return fmt.Errorf("game_servers[%d] (%s): port が不正です (%d)", i, s.Name, s.Port)
		}
		if s.TimeoutMS <= 0 {
			return fmt.Errorf("game_servers[%d] (%s): timeout_ms は1以上である必要があります", i, s.Name)
		}
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	fmt.Printf("[DEBUG] 設定ファイルを読み込みます: %s\n", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
