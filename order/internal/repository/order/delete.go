package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/order/internal/errors"
)

func (r *repository) Delete(ctx context.Context, orderUuid uuid.UUID) error {
	const deleteOrderQuery = "DELETE FROM orders WHERE uuid = $1"

	const deleteOrderItemsQuery = "DELETE from order_items WHERE order_uuid = $1"

	_, err := r.pool.Exec(ctx, deleteOrderQuery, orderUuid)

	if err != nil {
		return errs.ErrOrderNotFound
	}

	_, err = r.pool.Exec(ctx, deleteOrderItemsQuery, orderUuid)

	if err != nil {
		return errs.ErrOrderNotFound
	}

	return nil
}
