package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
)

func (r *repository) Pay(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) error {
	const query = "UPDATE orders SET payment_method = $1, transaction_uuid = $2 WHERE uuid = $3"

	_, err := r.pool.Exec(ctx, query, paymentMethod, transactionId, orderId)

	return err
}
