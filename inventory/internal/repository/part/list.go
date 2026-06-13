package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/1mpuser/inventory/internal/repository/record"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func (r *repository) ListPartsByUuids(ctx context.Context, uuids []uuid.UUID) ([]record.Part, error) {

	const query = "SELECT * from parts where uuid = ANY($1)"

	rows, err := r.pool.Query(ctx, query, uuids)

	if err != nil {
		return nil, err
	}

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])

	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (r *repository) ListPartsByPartType(ctx context.Context, partType inventoryv1.PartType) ([]record.Part, error) {

	const query = "SELECT * FROM parts where part_type = $1"

	rows, err := r.pool.Query(ctx, query, partType)

	if err != nil {
		return nil, err
	}

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])

	if err != nil {
		return nil, err
	}

	return parts, nil

}

func (r *repository) ListAllParts(ctx context.Context) ([]record.Part, error) {
	const query = "SELECT * FROM parts"

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])

}
