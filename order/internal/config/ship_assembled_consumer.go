package config

// ShipAssembledConsumerConfig — параметры consumer'а события ShipAssembled.
// OrderService слушает ShipAssembled (producer — AssemblyService) и переводит заказ в ASSEMBLED.
type ShipAssembledConsumerConfig struct {
	Topic   string `yaml:"topic" env:"SHIP_ASSEMBLED_TOPIC" env-default:"ship.assembled"`
	GroupID string `yaml:"group_id" env:"SHIP_ASSEMBLED_GROUP_ID" env-default:"order-service"`
}
