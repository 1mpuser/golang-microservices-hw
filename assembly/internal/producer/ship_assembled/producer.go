package shipassembled

import (
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/1mpuser/assembly/internal/model"
	"github.com/1mpuser/platform/pkg/kafka"
	"github.com/1mpuser/platform/pkg/kafka/producer"
	eventsv1 "github.com/1mpuser/shared/pkg/proto/events/v1"
)

// Producer публикует событие ShipAssembled в Kafka.
type Producer struct {
	producer *producer.Producer
}

func NewProducer(p *producer.Producer) *Producer {
	return &Producer{producer: p}
}

// ProduceShipAssembled сериализует событие в protobuf и отправляет в Kafka.
func (p *Producer) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	eventMessage := eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
		AssembledAt:  timestamppb.New(event.AssembledAt),
	}

	message, err := proto.Marshal(&eventMessage)
	if err != nil {
		return err
	}

	return p.producer.Send(ctx, &kafka.Message{
		Key:   []byte(event.OrderUUID),
		Value: message,
	})
}
