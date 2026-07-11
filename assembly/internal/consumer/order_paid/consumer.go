package orderpaid

import (
	"context"

	"github.com/1mpuser/platform/pkg/middleware/kafka"

	"github.com/1mpuser/platform/pkg/kafka/consumer"
	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer *consumer.Consumer
	handler  *Handler
}

func NewConsumer(group sarama.ConsumerGroup, topic string, handler *Handler) *Consumer {
	topicSlice := make([]string, 1)
	topicSlice[0] = topic

	return &Consumer{
		consumer: consumer.NewConsumer(group, topicSlice, consumer.WithMiddlewares(kafka.ConsumerLogging())),
		handler:  handler,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.handler.Handle)
}
