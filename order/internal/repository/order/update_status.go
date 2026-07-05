package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
)

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	conn := r.txGetter.DefaultTrOrDB(ctx, r.pool)

	const query = `
		UPDATE orders 
		SET status = $2
		WHERE uuid = $1
	`

	_, err := conn.Exec(ctx, query, id, status)

	return err
}
