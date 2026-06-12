package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	orderUUID, err := uuid.Parse(orderUuid)

	if err != nil {
		return err
	}

	err = s.orderRepository.Delete(ctx, orderUUID)

	return err
}
