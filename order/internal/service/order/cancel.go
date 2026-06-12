package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/order/internal/errors"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	orderUUID, err := uuid.Parse(orderUuid)
	if err != nil {
		return errs.ErrInvalidUUID
	}

	err = s.orderRepository.Delete(ctx, orderUUID)

	return err
}
