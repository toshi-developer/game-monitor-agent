package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Name      string `yaml:"name"`
	GameType  string `yaml:"game_type"`
	Address   string `yaml:"address"`
	Port      int    `yaml:"port"`
	TimeoutMS int    `yaml:"timeout_ms"`
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
