package session

// Реализовать (неделя 6): SessionRepository на Redis (go-redis v9, platform/pkg/redis).
//   Create: HSET session:{uuid} (поля из redisview) + EXPIRE (TTL 24h)
//   Get:    HGETALL session:{uuid}; пусто → errs.ErrSessionNotFound
//   Delete: DEL session:{uuid} (идемпотентно)
import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	errs "github.com/1mpuser/iam/internal/errors"
	redisview "github.com/1mpuser/iam/internal/repository/redis_view"
)

type repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *repository {
	return &repository{
		client,
	}
}

func (r *repository) Create(ctx context.Context, view redisview.SessionRedisView) error {
	key := "session:" + view.UUID

	pipe := r.client.TxPipeline()

	pipe.HSet(ctx, key, view)

	pipe.Expire(ctx, key, time.Duration(24)*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("сохранить сессию: %w", err)
	}

	return nil
}

func (r *repository) Get(ctx context.Context, uid string) (*redisview.SessionRedisView, error) {
	key := "session:" + uid

	res := r.client.HGetAll(ctx, key)

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("получить сессию: %w", err)
	}

	if len(res.Val()) == 0 {
		return nil, errs.ErrSessionNotFound
	}

	var session redisview.SessionRedisView

	if err := res.Scan(&session); err != nil {
		return nil, fmt.Errorf("распарсить сессию: %w", err)
	}

	return &session, nil
}

func (r *repository) Delete(ctx context.Context, uid string) error {
	key := "session:" + uid

	return r.client.Del(ctx, key).Err()
}
