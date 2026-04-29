// load config file to golang structs
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB        DBConfig        `yaml:"db"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Forecast  ForecastConfig  `yaml:"forecast"`
	Server    ServerConfig    `yaml:"server"`
}

type DBConfig struct {
	DSN string `yaml:"dsn"`
}

type SchedulerConfig struct {
	RunAt string `yaml:"run_at"`
}

type ForecastConfig struct {
	LookbackDays int `yaml:"lookback_days"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.Scheduler.RunAt == "" {
		cfg.Scheduler.RunAt = "0 2 * * *"
	}
	if cfg.Forecast.LookbackDays == 0 {
		cfg.Forecast.LookbackDays = 7
	}
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}

	return &cfg, nil
}
