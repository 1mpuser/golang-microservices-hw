package config

// AssemblerConfig — параметры эмуляции сборки корабля (случайная длительность в секундах).
type AssemblerConfig struct {
	MinBuildTimeSec int64 `yaml:"min_build_time_sec" env:"ASSEMBLER_MIN_BUILD_TIME_SEC" env-default:"5"`
	MaxBuildTimeSec int64 `yaml:"max_build_time_sec" env:"ASSEMBLER_MAX_BUILD_TIME_SEC" env-default:"15"`
}
