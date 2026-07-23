package config

// IAMClientConfig — адрес gRPC-сервера IAMService.
type IAMClientConfig struct {
	Address string `yaml:"address" env:"IAM_ADDRESS" env-default:"localhost:50053"`
}
