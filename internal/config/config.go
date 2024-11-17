package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

type LogLevel string

func (l LogLevel) Level() zap.AtomicLevel {
	switch l {
	case "info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "error":
		fallthrough
	default:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}
}

type Config struct {
	LogLevel LogLevel `yaml:"logLevel" env:"LOG_LEVEL,default=info"`
	ApiKey   string   `yaml:"apiKey" env:"API_KEY"`
	Endpoint string   `yaml:"endpoint" env:"API_ENDPOINT,default=https://driver-vehicle-licensing.api.gov.uk"`

	Interval time.Duration `yaml:"interval" env:"INTERVAL,default=1h"`
	Port     int           `yaml:"port" env:"PORT,default=9876"`

	Vehicles []string `yaml:"vehicles" env:"vehicles"`
}

func Load(path string) (*Config, error) {
	conf := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(file, conf); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

	if err := envconfig.Process(context.Background(), conf); err != nil {
		return nil, fmt.Errorf("failed to process env vars: %w", err)
	}
	return conf, nil
}
