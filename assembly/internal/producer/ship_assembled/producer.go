package shipassembled

import (
	"context"

	"github.com/1mpuser/assembly/internal/model"
)

// Producer публикует событие ShipAssembled в Kafka.
//
// TODO(неделя 5, часть 4): реализовать —
//  - хранить *producer.Producer (platform) привязанный к топику ship.assembled;
//  - конструктор NewProducer(p *producer.Producer) *Producer.
type Producer struct{}

// ProduceShipAssembled сериализует событие в protobuf и отправляет в Kafka.
//
// TODO(неделя 5): смаппить model.ShipAssembled → eventsv1.ShipAssembled, proto.Marshal,
// собрать kafka.Message (Key = order_uuid для сохранения порядка по заказу) и отправить через Send.
func (p *Producer) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	panic("не реализовано: часть 4 (producer ShipAssembled)")
}
