package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, uuid uuid.UUID) (record.Part, error) {

	const query = "SELECT * from parts where id = $1"

	row, err := r.pool.Query(ctx, query, uuid)

	if err != nil {
		return record.Part{}, errs.ErrPartNotFound
	}

	part, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[record.Part])

	if err != nil {
		return record.Part{}, errs.ErrPartNotFound
	}

	return part, nil
}
