package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		HTTPPort string `yaml:"http_port"`
		GRPCPort string `yaml:"grpc_port"`
	} `yaml:"server"`
	Storage struct {
		Type     string `yaml:"type"`
		Postgres struct {
			URL string `yaml:"url"`
		} `yaml:"postgres"`
	} `yaml:"storage"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("не удалось декодировать конфигурацию: %w", err)
	}

	// Установка значений по умолчанию
	if config.Server.HTTPPort == "" {
		config.Server.HTTPPort = "8080"
	}
	if config.Server.GRPCPort == "" {
		config.Server.GRPCPort = "50051"
	}
	if config.Storage.Type == "" {
		config.Storage.Type = "inmemory"
	}

	return config, nil
}
