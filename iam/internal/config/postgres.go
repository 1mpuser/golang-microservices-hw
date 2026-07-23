package config

import "fmt"

// PGConfig — параметры подключения к PostgreSQL.
type PGConfig struct {
	Host     string `yaml:"host"     env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     string `yaml:"port"     env:"POSTGRES_PORT"     env-default:"7434"`
	Database string `yaml:"database" env:"POSTGRES_DB"       env-default:"iam-service"`
	User     string `yaml:"user"     env:"POSTGRES_USER"     env-default:"iam-service-user"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"iam-service-password"`
	SSLMode  string `yaml:"sslmode"  env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

// DSN — строка подключения для pgxpool.New.
func (c *PGConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Database, c.User, c.Password, c.SSLMode,
	)
}
