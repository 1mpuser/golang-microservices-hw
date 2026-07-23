package user

import (
	"context"

	"github.com/jackc/pgx/v5"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/repository/record"
)

func (r *repository) GetByLogin(ctx context.Context, login string) (*record.User, error) {
	const query = `
		SELECT * from users
		WHERE login = $1
	`
	rows, err := r.pool.Query(ctx, query, login)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[record.User])
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}
