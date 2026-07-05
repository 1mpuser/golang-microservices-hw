package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config — корневая конфигурация InventoryService.
type Config struct {
	GRPC   GRPCConfig   `yaml:"grpc"`
	PG     PGConfig     `yaml:"pg"`
	Logger LoggerConfig `yaml:"logger"`
}

var appConfig *Config

// MustLoad загружает конфигурацию или завершает процесс с паникой.
func MustLoad() {
	cfg, err := Load(ResolveConfigPath())
	if err != nil {
		panic(fmt.Sprintf("не удалось загрузить конфиг: %v", err))
	}
	appConfig = cfg
}

// AppConfig возвращает загруженный конфиг (должен быть вызван после MustLoad).
func AppConfig() *Config {
	return appConfig
}

const defaultConfigPath = "config.local.yaml"

// ResolveConfigPath определяет путь к конфиг-файлу по цепочке приоритетов:
// флаг -config > env CONFIG_PATH > "config.local.yaml".
func ResolveConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}
	return defaultConfigPath
}

// Load загружает конфигурацию из YAML-файла и env-переменных поверх.
func Load(path string) (*Config, error) {
	var cfg Config

	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("загрузить конфиг из %q: %w", path, err)
		}
		return &cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("загрузить конфиг из env: %w", err)
	}
	return &cfg, nil
}
