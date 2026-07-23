package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/repository/record"
)

func (r *repository) Create(ctx context.Context, user record.User) error {
	const insertOrderQuery = `
		INSERT INTO users 
		(uuid, login, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(
		ctx, insertOrderQuery,
		user.UUID,
		user.Login,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrUserAlreadyExists
		}
		return fmt.Errorf("создать пользователя: %w", err)
	}

	return nil
}
