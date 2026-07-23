package iam

import (
	"context"

	"github.com/1mpuser/iam/internal/repository/record"
	redisview "github.com/1mpuser/iam/internal/repository/redis_view"
)

// UserRepository — то, что сервису нужно от хранилища пользователей (Postgres).
// Интерфейс объявлен здесь, у потребителя: сервис зависит от абстракции,
// а конкретный repository/user её реализует. Отсюда же mockery берёт его для тестов.
type UserRepository interface {
	Create(ctx context.Context, user record.User) error
	GetByLogin(ctx context.Context, login string) (*record.User, error)
	GetByUuid(ctx context.Context, uid string) (*record.User, error)
}

// SessionRepository — то, что сервису нужно от хранилища сессий (Redis).
type SessionRepository interface {
	Create(ctx context.Context, view redisview.SessionRedisView) error
	Get(ctx context.Context, uid string) (*redisview.SessionRedisView, error)
	Delete(ctx context.Context, uid string) error
}
