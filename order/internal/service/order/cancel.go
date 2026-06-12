package order

import (
	"context"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/google/uuid"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	orderUUID, err := uuid.Parse(orderUuid)

	if err != nil {
		return errs.ErrInvalidUUID
	}

	err = s.orderRepository.Delete(ctx, orderUUID)

	return err
}
