package config

// OrderPaidConsumerConfig — параметры consumer'а события OrderPaid.
type OrderPaidConsumerConfig struct {
	Topic   string `yaml:"topic" env:"ORDER_PAID_TOPIC" env-default:"order.paid"`
	GroupID string `yaml:"group_id" env:"ORDER_PAID_GROUP_ID" env-default:"assembly-service"`
}
