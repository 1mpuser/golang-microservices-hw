// Package app — тонкая обёртка над сборкой IAMService для API-тестов.
// Собирает зависимости так же, как internal/app/di.go, но принимает пул БД и
// Redis-клиент напрямую (без чтения конфига) — для запуска на bufconn.
package app

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	authAPIv1 "github.com/1mpuser/iam/internal/api/auth/v1"
	userAPIv1 "github.com/1mpuser/iam/internal/api/user/v1"
	"github.com/1mpuser/iam/internal/interceptor"
	sessionrepo "github.com/1mpuser/iam/internal/repository/session"
	userrepo "github.com/1mpuser/iam/internal/repository/user"
	iamservice "github.com/1mpuser/iam/internal/service/iam"
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
	userv1 "github.com/1mpuser/shared/pkg/proto/user/v1"
)

// NewGRPCServer собирает IAMService на переданных pool/rdb и возвращает готовый
// gRPC-сервер с error-interceptor и зарегистрированными AuthService + UserService.
//
// sessionTTL и bcryptCost приходят из теста (короткий TTL, MinCost для скорости).
// Текущая реализация сервиса зашивает TTL/cost внутри — параметры зарезервированы
// (чтобы задействовать их, нужно протянуть значения через NewService/login/register).
func NewGRPCServer(pool *pgxpool.Pool, rdb *goredis.Client, _ time.Duration, _ int) *grpc.Server {
	userRepo := userrepo.New(pool)
	sessionRepo := sessionrepo.NewRepository(rdb)
	svc := iamservice.NewService(userRepo, sessionRepo)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ErrorUnaryInterceptor))

	authv1.RegisterAuthServiceServer(server, authAPIv1.NewApi(svc))
	userv1.RegisterUserServiceServer(server, userAPIv1.NewApi(svc))

	return server
}
