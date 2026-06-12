package order

import (
	"context"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Pay(_ context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[orderId]

	if !ok {
		return errs.ErrPartNotFound
	}

	order.Status = model.OrderStatusPaid

	order.PaymentMethod = &paymentMethod

	order.TransactionUUID = new(transactionId)

	r.data[orderId] = order

	return nil
}
