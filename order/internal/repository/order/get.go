package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (*record.Order, error) {
	conn := r.txGetter.DefaultTrOrDB(ctx, r.pool)

	const query = `
		SELECT
		o.uuid,
		o.total_price,
		o.status,
		o.transaction_uuid,
		o.payment_method,
		o.created_at,
		MAX(CASE WHEN oi.part_type = 'HULL'   THEN oi.part_uuid::text END)::uuid AS hull_uuid,
		MAX(CASE WHEN oi.part_type = 'ENGINE' THEN oi.part_uuid::text END)::uuid AS engine_uuid,
		MAX(CASE WHEN oi.part_type = 'SHIELD' THEN oi.part_uuid::text END)::uuid AS shield_uuid,
		MAX(CASE WHEN oi.part_type = 'WEAPON' THEN oi.part_uuid::text END)::uuid AS weapon_uuid
		FROM orders o
		LEFT JOIN order_items oi ON oi.order_uuid = o.uuid
		WHERE o.uuid = $1
		GROUP BY o.uuid, o.total_price, o.status, o.transaction_uuid, o.payment_method, o.created_at
	`

	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, errs.ErrOrderNotFound
	}

	order, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[record.Order])
	if err != nil {
		return nil, errs.ErrOrderNotFound
	}

	return order, nil
}
