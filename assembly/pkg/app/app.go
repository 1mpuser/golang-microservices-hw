// Package app — тонкая обёртка над internal-пакетами AssemblyService для e2e-тестов:
// собирает всю цепочку (consumer OrderPaid → service → producer ShipAssembled)
// из переданных sarama-продюсера и consumer-группы.
package app

import (
	"context"

	"github.com/IBM/sarama"

	orderPaidConsumer "github.com/1mpuser/assembly/internal/consumer/order_paid"
	shipAssembledProducer "github.com/1mpuser/assembly/internal/producer/ship_assembled"
	assemblyService "github.com/1mpuser/assembly/internal/service/assembly"
	kafkaProducer "github.com/1mpuser/platform/pkg/kafka/producer"
)

// Config — параметры запуска AssemblyService.
type Config struct {
	OrderPaidTopic     string
	ShipAssembledTopic string
	MinBuildTimeSec    int64
	MaxBuildTimeSec    int64
}

// App — собранный AssemblyService, готовый к запуску consumer'а.
type App struct {
	consumer *orderPaidConsumer.Consumer
}

// New собирает AssemblyService из готовых sarama sync-продюсера, consumer-группы и конфига.
func New(syncProducer sarama.SyncProducer, group sarama.ConsumerGroup, cfg Config) *App {
	platformProducer := kafkaProducer.NewProducer(syncProducer, cfg.ShipAssembledTopic)
	shipProducer := shipAssembledProducer.NewProducer(platformProducer)

	svc := assemblyService.NewService(shipProducer, cfg.MinBuildTimeSec, cfg.MaxBuildTimeSec)
	handler := orderPaidConsumer.NewHandler(svc)
	consumer := orderPaidConsumer.NewConsumer(group, cfg.OrderPaidTopic, handler)

	return &App{consumer: consumer}
}

// RunConsumer запускает чтение OrderPaid до отмены ctx.
func (a *App) RunConsumer(ctx context.Context) error {
	return a.consumer.Run(ctx)
}
