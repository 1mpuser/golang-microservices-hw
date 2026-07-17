package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	orderUUID, err := uuid.Parse(orderUuid)
	if err != nil {
		return errs.ErrInvalidUUID
	}

	return s.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.orderRepository.Get(ctx, orderUUID)
		if err != nil {
			return err
		}

		switch order.Status {
		case model.OrderStatusPaid:
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			return errs.ErrOrderCancelled
		case model.OrderStatusAssembled:
			return errs.ErrOrderAssembled
		}

		partIds := make([]string, 0)

		partIds = append(partIds, order.HullUUID.String())
		partIds = append(partIds, order.EngineUUID.String())

		if order.ShieldUUID != nil {
			partIds = append(partIds, order.ShieldUUID.String())
		}

		if order.WeaponUUID != nil {
			partIds = append(partIds, order.WeaponUUID.String())
		}

		if err := s.inventoryClient.ReleaseParts(ctx, partIds); err != nil {
			return err
		}

		return s.orderRepository.UpdateStatus(ctx, orderUUID, model.OrderStatusCancelled)
	})
}
