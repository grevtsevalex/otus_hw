package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

// Config модель конфига.
type Config struct {
	Logger LoggerConf
	Queue  QueueConf
	Cron   Cron
}

// LoggerConf модель конфига логгера.
type LoggerConf struct {
	Level string `toml:"level"`
}

// StorageConf модель конфига очереди.
type QueueConf struct {
	Address string `toml:"address"`
}

// Cron модель конфига крона обработки очереди событий.
type Cron struct {
	Period int `toml:"periodSec"`
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
