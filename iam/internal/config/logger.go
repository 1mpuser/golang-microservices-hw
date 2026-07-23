package config

// LoggerConfig — параметры логгера.
type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}
