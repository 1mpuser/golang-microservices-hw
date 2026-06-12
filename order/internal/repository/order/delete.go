package order

import (
	"context"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Delete(_ context.Context, orderUuid uuid.UUID) error {
	r.mu.RLock()

	order, ok := r.data[orderUuid]

	if !ok {
		r.mu.RUnlock()
		return errs.ErrOrderNotFound
	}

	r.mu.RUnlock()

	switch order.Status {
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderCancelled
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	order.Status = model.OrderStatusCancelled

	r.data[order.OrderUUID] = order

	return nil
}
