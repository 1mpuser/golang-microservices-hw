package config

// ShipAssembledProducerConfig — параметры producer'а события ShipAssembled.
type ShipAssembledProducerConfig struct {
	Topic string `yaml:"topic" env:"SHIP_ASSEMBLED_TOPIC" env-default:"ship.assembled"`
}
