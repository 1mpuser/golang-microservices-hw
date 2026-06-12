package order

import (
	"context"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, orderId string) (model.Order, error) {
	orderUuid, err := uuid.Parse(orderId)

	if err != nil {
		return model.Order{}, errs.ErrInvalidUUID
	}

	order, err := s.orderRepository.Get(ctx, orderUuid)

	if err != nil {
		return model.Order{}, err
	}

	return converter.RecordToModel(order), nil
}
