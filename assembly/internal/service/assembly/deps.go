package assembly

import (
	"context"

	"github.com/1mpuser/assembly/internal/model"
)

// ShipAssembledProducer публикует событие о завершении сборки в Kafka.
type ShipAssembledProducer interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}
