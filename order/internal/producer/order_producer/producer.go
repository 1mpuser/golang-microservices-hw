package orderproducer

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/platform/pkg/kafka"
	"github.com/1mpuser/platform/pkg/kafka/producer"
	eventsv1 "github.com/1mpuser/shared/pkg/proto/events/v1"
)

type Producer struct {
	producer *producer.Producer
}

func New(p *producer.Producer) *Producer {
	return &Producer{producer: p}
}

func (p *Producer) ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error {
	eventMessage := eventsv1.OrderPaid{
		EventUuid: event.EventUUID,
		OrderUuid: event.OrderUUID,
		UserUuid:  event.UserUUID,
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
