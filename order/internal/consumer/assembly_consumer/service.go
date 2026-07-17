package assemblyconsumer

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/platform/pkg/kafka"
	"github.com/1mpuser/platform/pkg/kafka/consumer"
)

// Service слушает топик ShipAssembled и завершает сборку заказа:
// списывает детали (CommitParts) и переводит заказ в статус ASSEMBLED.
type Service struct {
	consumer  *consumer.Consumer
	repo      OrderRepository
	invClient InventoryClient
	txManager TxManager
}

// NewService собирает consumer ShipAssembled с его зависимостями.
func NewService(c *consumer.Consumer, repo OrderRepository, invClient InventoryClient, txManager TxManager) *Service {
	return &Service{
		consumer:  c,
		repo:      repo,
		invClient: invClient,
		txManager: txManager,
	}
}

// RunConsumer запускает бесконечное чтение ShipAssembled до отмены ctx.
func (s *Service) RunConsumer(ctx context.Context) error {
	return s.consumer.Consume(ctx, s.handle)
}

// handle обрабатывает одно событие ShipAssembled.
//
// Идемпотентность обеспечивается статусом заказа в БД: повторное событие для уже
// собранного заказа ничего не делает (at-least-once, переживает рестарт).
func (s *Service) handle(ctx context.Context, msg kafka.Message) error {
	event, err := decodeShipAssembled(msg.Value)
	if err != nil {
		return err
	}

	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		return err
	}

	return s.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.repo.GetForUpdate(ctx, orderUUID)
		if err != nil {
			return err
		}

		// Уже собран — дубликат события, тихо пропускаем.
		if order.Status == model.OrderStatusAssembled {
			return nil
		}

		// Собирать можно только оплаченный заказ; прочие статусы игнорируем.
		if order.Status != model.OrderStatusPaid {
			return nil
		}

		partIDs := make([]string, 0, 4)
		partIDs = append(partIDs, order.HullUUID.String(), order.EngineUUID.String())
		if order.ShieldUUID != nil {
			partIDs = append(partIDs, order.ShieldUUID.String())
		}
		if order.WeaponUUID != nil {
			partIDs = append(partIDs, order.WeaponUUID.String())
		}

		if err := s.invClient.CommitParts(ctx, partIDs); err != nil {
			return err
		}

		return s.repo.UpdateStatus(ctx, orderUUID, model.OrderStatusAssembled)
	})
}
