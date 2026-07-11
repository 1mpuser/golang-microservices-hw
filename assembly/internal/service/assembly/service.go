package assembly

// service — application-сервис сборки корабля.
// TODO(неделя 5): при необходимости добавить зависимости (например, параметры
// длительности сборки из config.AssemblerConfig).
type service struct {
	producer ShipAssembledProducer
}

// NewService создаёт сервис сборки.
func NewService(producer ShipAssembledProducer) *service {
	return &service{
		producer: producer,
	}
}
