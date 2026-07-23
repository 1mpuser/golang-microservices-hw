package user

import (
	"context"

	"github.com/jackc/pgx/v5"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/repository/record"
)

func (r *repository) GetByUuid(ctx context.Context, uid string) (*record.User, error) {
	const query = `
		SELECT * from users
		WHERE uuid = $1
	`

	rows, err := r.pool.Query(ctx, query, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[record.User])
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}
