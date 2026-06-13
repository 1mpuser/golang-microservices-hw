package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/repository/record"
)

func (r *repository) Create(ctx context.Context, order record.Order, orderItems []record.OrderItem) error {
	conn := r.txGetter.DefaultTrOrDB(ctx, r.pool)

	const insertOrderQuery = `
		INSERT INTO orders (uuid, total_price, status, transaction_uuid, payment_method, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := conn.Exec(ctx, insertOrderQuery,
		order.OrderUUID,
		order.TotalPrice,
		order.Status,
		order.TransactionUUID,
		order.PaymentMethod,
		order.CreatedAt,
	)
	if err != nil {
		return err
	}

	orderUUIDs := make([]uuid.UUID, 0, len(orderItems))
	partUUIDs := make([]uuid.UUID, 0, len(orderItems))
	partTypes := make([]string, 0, len(orderItems))
	prices := make([]int64, 0, len(orderItems))
	createdAts := make([]time.Time, 0, len(orderItems))

	for _, item := range orderItems {
		orderUUIDs = append(orderUUIDs, item.OrderUUID)
		partUUIDs = append(partUUIDs, item.PartUUID)
		partTypes = append(partTypes, string(item.PartType))
		prices = append(prices, item.Price)
		createdAts = append(createdAts, item.CreatedAt)
	}

	const insertItemsQuery = `
		INSERT INTO order_items (order_uuid, part_uuid, part_type, price, created_at)
		SELECT * FROM unnest($1::uuid[], $2::uuid[], $3::text[], $4::bigint[], $5::timestamp[])
	`

	_, err = conn.Exec(ctx, insertItemsQuery, orderUUIDs, partUUIDs, partTypes, prices, createdAts)

	return err
}
