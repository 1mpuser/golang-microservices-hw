package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	authAPIv1 "github.com/1mpuser/iam/internal/api/auth/v1"
	userAPIv1 "github.com/1mpuser/iam/internal/api/user/v1"
	"github.com/1mpuser/iam/internal/config"
	sessionrepo "github.com/1mpuser/iam/internal/repository/session"
	userrepo "github.com/1mpuser/iam/internal/repository/user"
	iamservice "github.com/1mpuser/iam/internal/service/iam"
	"github.com/1mpuser/platform/pkg/closer"
	platformredis "github.com/1mpuser/platform/pkg/redis" // ← фабрика NewClient
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
	userv1 "github.com/1mpuser/shared/pkg/proto/user/v1"
)

// Реализовать (неделя 6): ручной DI-контейнер (без wire/fx).
// pgxpool → Redis-клиент (platform/pkg/redis) → repository(user, session) →
// service(iam) → api(auth, user) → grpc.NewServer(error-interceptor из internal/interceptor).

type iamFacade interface {
	authAPIv1.IAmService // Login, Whoami, Logout
	userAPIv1.IAmService // Register, GetUser
}

type diContainer struct {
	pgPool      *pgxpool.Pool
	redisClient *goredis.Client
	userRepo    iamservice.UserRepository
	sessionRepo iamservice.SessionRepository
	iamSvc      iamFacade
	authAPI     authv1.AuthServiceServer
	userAPI     userv1.UserServiceServer
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if nil == d.pgPool {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("Postgres", "error", err)
			os.Exit(1)
		}
		if err = pool.Ping(ctx); err != nil {
			slog.Error("Postgres ping", "error", err)
			os.Exit(1)
		}

		closer.Add("Postgres pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) RedisClient(ctx context.Context) *goredis.Client {
	if d.redisClient == nil {
		cli, err := platformredis.NewClient(
			&goredis.Options{
				Addr:     config.AppConfig().Redis.Address(), // "localhost:6379"
				Password: config.AppConfig().Redis.Password,
				DB:       config.AppConfig().Redis.DB,
			}, slog.Default(),
		)
		if err != nil {
			slog.Error("не удалось подключиться к Redis", "error", err)
			os.Exit(1)
		}

		closer.Add("Redis", func(_ context.Context) error {
			return cli.Close()
		})

		d.redisClient = cli
	}

	return d.redisClient
}

func (d *diContainer) UserRepository(ctx context.Context) iamservice.UserRepository {
	if nil == d.userRepo {
		d.userRepo = userrepo.New(d.pgPool)
	}

	return d.userRepo
}

func (d *diContainer) SessionRepository(ctx context.Context) iamservice.SessionRepository {
	if nil == d.sessionRepo {
		d.sessionRepo = sessionrepo.NewRepository(d.RedisClient(ctx))
	}

	return d.sessionRepo
}

func (d *diContainer) IAmService(ctx context.Context) iamFacade {
	if nil == d.iamSvc {
		d.iamSvc = iamservice.NewService(d.UserRepository(ctx), d.SessionRepository(ctx))
	}

	return d.iamSvc
}

func (d *diContainer) AuthAPI(ctx context.Context) authv1.AuthServiceServer {
	if d.authAPI == nil {
		d.authAPI = authAPIv1.NewApi(d.IAmService(ctx))
	}
	return d.authAPI
}

func (d *diContainer) UserAPI(ctx context.Context) userv1.UserServiceServer {
	if d.userAPI == nil {
		d.userAPI = userAPIv1.NewApi(d.IAmService(ctx))
	}
	return d.userAPI
}
