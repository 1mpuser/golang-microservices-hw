package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/record"
)

func (r *repository) ListPartsByUuids(ctx context.Context, uuids []uuid.UUID) ([]record.PartRecord, error) {
	const query = "SELECT * from parts where uuid = ANY($1) ORDER BY array_position($1, uuid)"

	rows, err := r.pool.Query(ctx, query, uuids)
	if err != nil {
		return nil, err
	}

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.PartRecord])
	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (r *repository) ListPartsByPartType(ctx context.Context, partType model.PartType) ([]record.PartRecord, error) {
	const query = "SELECT * FROM parts where part_type = $1"

	rows, err := r.pool.Query(ctx, query, partType)
	if err != nil {
		return nil, err
	}

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.PartRecord])
	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (r *repository) ListAllParts(ctx context.Context) ([]record.PartRecord, error) {
	const query = "SELECT * FROM parts"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[record.PartRecord])
}

func (r *repository) ListPartsForUpdate(ctx context.Context, uuids []uuid.UUID) ([]record.PartRecord, error) {
	// Берём транзакционное соединение — иначе FOR UPDATE залочит строки на отдельном
	// коннекте пула и тут же отпустит их (autocommit), и блокировка не сработает.
	conn := r.txGetter.DefaultTrOrDB(ctx, r.pool)

	const query = "SELECT * from parts where uuid = ANY($1) ORDER BY array_position($1, uuid) FOR UPDATE"

	rows, err := conn.Query(ctx, query, uuids)
	if err != nil {
		return nil, err
	}

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.PartRecord])
	if err != nil {
		return nil, err
	}

	return parts, nil
}
