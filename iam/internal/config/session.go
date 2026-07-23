package config

import "time"

// SessionConfig — параметры пользовательских сессий.
type SessionConfig struct {
	TTL time.Duration `yaml:"ttl" env:"SESSION_TTL" env-default:"24h"`
}
