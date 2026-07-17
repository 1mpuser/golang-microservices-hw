package assembly

// service — application-сервис сборки корабля.
type service struct {
	producer        ShipAssembledProducer
	minBuildSeconds int64
	maxBuildSeconds int64
}

// NewService создаёт сервис сборки. Диапазон [minBuildSeconds, maxBuildSeconds]
// задаёт эмулируемое время сборки (0/0 — мгновенно, для e2e-тестов).
func NewService(producer ShipAssembledProducer, minBuildSeconds, maxBuildSeconds int64) *service {
	return &service{
		producer:        producer,
		minBuildSeconds: minBuildSeconds,
		maxBuildSeconds: maxBuildSeconds,
	}
}
