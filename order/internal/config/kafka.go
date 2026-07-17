package config

// KafkaConfig — общие параметры подключения к Kafka.
type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
}
