package config

import "net"

// GRPCConfig — параметры gRPC-сервера.
type GRPCConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
}

func (c *GRPCConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
