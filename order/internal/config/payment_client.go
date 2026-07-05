package config

// PaymentClientConfig — адрес gRPC-сервера PaymentService.
type PaymentClientConfig struct {
	Address string `yaml:"address" env:"PAYMENT_ADDRESS" env-default:"localhost:50052"`
}
