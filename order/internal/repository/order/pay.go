package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
)

func (r *repository) Pay(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) error {
	conn := r.txGetter.DefaultTrOrDB(ctx, r.pool)

	const query = `
		UPDATE orders
		SET
			payment_method = $1,
			transaction_uuid = $2,
			status = $3
		WHERE uuid = $4
	`

	_, err := conn.Exec(ctx, query, paymentMethod, transactionId, model.OrderStatusPaid, orderId)

	return err
}
