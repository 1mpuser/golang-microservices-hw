package app

// diContainer — контейнер зависимостей AssemblyService (ленивая инициализация через геттеры).
//
// TODO(неделя 5, часть 4): добавить поля и геттеры с nil-check —
//   - sarama.SyncProducer и sarama.ConsumerGroup из config.Kafka (+ регистрация в closer);
//   - producer ShipAssembled (обёртка над platform producer);
//   - service Assembly (использует producer);
//   - handler OrderPaid (использует service);
//   - consumer OrderPaid (group + topic + handler, middleware логирования).
type diContainer struct{}

func newDiContainer() *diContainer {
	return &diContainer{}
}
