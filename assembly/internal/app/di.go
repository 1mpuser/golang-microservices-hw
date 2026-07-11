package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/IBM/sarama"

	"github.com/1mpuser/assembly/internal/config"
	orderPaidConsumer "github.com/1mpuser/assembly/internal/consumer/order_paid"
	shipAssembledProducer "github.com/1mpuser/assembly/internal/producer/ship_assembled"
	assemblyService "github.com/1mpuser/assembly/internal/service/assembly"
	"github.com/1mpuser/platform/pkg/closer"
	"github.com/1mpuser/platform/pkg/kafka/producer"
)

// diContainer — контейнер зависимостей AssemblyService с ленивой инициализацией через геттеры.
type diContainer struct {
	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
	shipProducer  *shipAssembledProducer.Producer
	service       orderPaidConsumer.AssemblyService
	handler       *orderPaidConsumer.Handler
	consumer      *orderPaidConsumer.Consumer
}

func newDiContainer() *diContainer {
	return &diContainer{}
}

// SyncProducer — sarama SyncProducer для отправки событий ShipAssembled.
func (d *diContainer) SyncProducer(_ context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		// Обязательно для SyncProducer — иначе SendMessage зависнет навсегда.
		cfg.Producer.Return.Successes = true
		// Ждём подтверждения записи от всех in-sync реплик (надёжность).
		cfg.Producer.RequiredAcks = sarama.WaitForAll

		sp, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers, cfg)
		if err != nil {
			slog.Error("не удалось создать Kafka sync producer", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka sync producer", func(_ context.Context) error {
			return sp.Close()
		})

		d.syncProducer = sp
	}

	return d.syncProducer
}

// ShipAssembledProducer — доменный продюсер события ShipAssembled (обёртка над platform producer).
func (d *diContainer) ShipAssembledProducer(ctx context.Context) *shipAssembledProducer.Producer {
	if d.shipProducer == nil {
		platformProducer := producer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig().ShipAssembledProducer.Topic,
		)
		d.shipProducer = shipAssembledProducer.NewProducer(platformProducer)
	}

	return d.shipProducer
}

// AssemblyService — сервис сборки корабля.
func (d *diContainer) AssemblyService(ctx context.Context) orderPaidConsumer.AssemblyService {
	if d.service == nil {
		d.service = assemblyService.NewService(d.ShipAssembledProducer(ctx))
	}

	return d.service
}

// OrderPaidHandler — обработчик входящего события OrderPaid.
func (d *diContainer) OrderPaidHandler(ctx context.Context) *orderPaidConsumer.Handler {
	if d.handler == nil {
		d.handler = orderPaidConsumer.NewHandler(d.AssemblyService(ctx))
	}

	return d.handler
}

// ConsumerGroup — sarama ConsumerGroup для чтения топика OrderPaid.
func (d *diContainer) ConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		// Читать с начала топика, если для группы ещё нет закоммиченного оффсета,
		// иначе новый consumer group пропустит уже лежащие в топике события.
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().OrderPaidConsumer.GroupID,
			cfg,
		)
		if err != nil {
			slog.Error("не удалось создать Kafka consumer group", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}

// OrderPaidConsumer — consumer события OrderPaid (верхушка графа зависимостей).
func (d *diContainer) OrderPaidConsumer(ctx context.Context) *orderPaidConsumer.Consumer {
	if d.consumer == nil {
		d.consumer = orderPaidConsumer.NewConsumer(
			d.ConsumerGroup(ctx),
			config.AppConfig().OrderPaidConsumer.Topic,
			d.OrderPaidHandler(ctx),
		)
	}

	return d.consumer
}
