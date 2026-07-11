package config

import "time"

// AssemblerConfig — параметры эмуляции сборки корабля (случайная длительность в диапазоне).
type AssemblerConfig struct {
	MinBuildTime time.Duration `yaml:"min_build_time" env:"ASSEMBLER_MIN_BUILD_TIME" env-default:"5s"`
	MaxBuildTime time.Duration `yaml:"max_build_time" env:"ASSEMBLER_MAX_BUILD_TIME" env-default:"15s"`
}
