package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/repository/converter"
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
