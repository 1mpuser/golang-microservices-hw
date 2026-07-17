package config

// OrderPaidProducerConfig — параметры producer'а события OrderPaid.
// OrderService публикует OrderPaid после успешной оплаты; consumer — AssemblyService.
type OrderPaidProducerConfig struct {
	Topic string `yaml:"topic" env:"ORDER_PAID_TOPIC" env-default:"order.paid"`
}
