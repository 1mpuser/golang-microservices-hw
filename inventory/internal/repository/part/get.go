package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, uuid uuid.UUID) (record.PartRecord, error) {
	const query = "SELECT * from parts where uuid = $1"

	row, err := r.pool.Query(ctx, query, uuid)
	if err != nil {
		return record.PartRecord{}, errs.ErrPartNotFound
	}

	part, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[record.PartRecord])
	if err != nil {
		return record.PartRecord{}, errs.ErrPartNotFound
	}

	return part, nil
}
