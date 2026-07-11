package part

import (
	"context"

	"github.com/1mpuser/inventory/internal/model"
)

func (r *repository) UpdateParts(ctx context.Context, parts []model.Part) error {
	if len(parts) == 0 {
		return nil
	}

	conn := r.txGetter.DefaultTrOrDB(ctx, r.pool)

	const query = `
		UPDATE parts AS p
		SET
			name           = t.name,
			description    = t.description,
			price          = t.price,
			stock_quantity = t.stock_quantity,
			reserved       = t.reserved
		FROM unnest(
			$1::uuid[],
			$2::text[],
			$3::text[],
			$4::bigint[],
			$5::int[],
			$6::int[]
		) AS t(uuid, name, description, price, stock_quantity, reserved)
		WHERE p.uuid = t.uuid`

	// параллельные массивы: по одному на каждую колонку из unnest
	uuids := make([]string, len(parts))
	names := make([]string, len(parts))
	descriptions := make([]string, len(parts))
	prices := make([]int64, len(parts))
	stockQuantities := make([]int, len(parts))
	reserved := make([]int, len(parts))

	for i := range parts {
		uuids[i] = parts[i].UUID()
		names[i] = parts[i].Name()
		descriptions[i] = parts[i].Description()
		prices[i] = parts[i].Price()
		stockQuantities[i] = parts[i].StockQuantity()
		reserved[i] = parts[i].Reserved()
	}

	_, err := conn.Exec(
		ctx,
		query,
		uuids,
		names,
		descriptions,
		prices,
		stockQuantities,
		reserved,
	)

	return err
}
