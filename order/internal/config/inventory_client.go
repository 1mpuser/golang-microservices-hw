package config

// InventoryClientConfig — адрес gRPC-сервера InventoryService.
type InventoryClientConfig struct {
	Address string `yaml:"address" env:"INVENTORY_ADDRESS" env-default:"localhost:50051"`
}
