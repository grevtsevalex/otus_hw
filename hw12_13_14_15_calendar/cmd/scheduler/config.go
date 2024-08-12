package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

// Config модель конфига.
type Config struct {
	Logger  LoggerConf
	DB      DBConf
	Storage StorageConf
	Queue   QueueConf
	Cron    Cron
}

// LoggerConf модель конфига логгера.
type LoggerConf struct {
	Level string `toml:"level"`
}

// StorageConf модель конфига хранилища.
type StorageConf struct {
	Type string `toml:"type"`
}

// StorageConf модель конфига очереди.
type QueueConf struct {
	Address string `toml:"address"`
}

// Cron модель конфига крона удаления событий.
type Cron struct {
	Period int `toml:"periodSec"`
}

// DBConf модель конфига БД.
type DBConf struct {
	Name            string `toml:"name"`
	User            string `toml:"user"`
	Pass            string `toml:"pass"`
	PoolSize        int    `toml:"poolSize"`
	MaxConnLifeTime int    `toml:"maxConnLifeTime"`
}

// NewConfig инициализация конфига.
func NewConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("reading config file: %w", err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		return config, fmt.Errorf("reading config file: %w", err)
	}

	return config, nil
}
