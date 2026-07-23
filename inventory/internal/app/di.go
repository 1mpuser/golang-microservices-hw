package app

import (
	"context"
	"log/slog"
	"os"

	iamGRPCClient "github.com/1mpuser/inventory/internal/client/grpc/iam/v1"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryAPIv1 "github.com/1mpuser/inventory/internal/api/inventory/v1"
	"github.com/1mpuser/inventory/internal/config"
	"github.com/1mpuser/inventory/internal/interceptor"
	partRepository "github.com/1mpuser/inventory/internal/repository/part"
	applicationPart "github.com/1mpuser/inventory/internal/service/application/part"
	domainService "github.com/1mpuser/inventory/internal/service/domain"
	"github.com/1mpuser/platform/pkg/closer"
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	pgPool           *pgxpool.Pool
	partRepo         applicationPart.PartRepository
	partSvc          inventoryAPIv1.PartService
	inventoryHandler inventoryv1.InventoryServiceServer
	txManager        applicationPart.TxManager
	iamCONN          *grpc.ClientConn
	iamClient        interceptor.IAMClient
	authInterceptor  grpc.UnaryServerInterceptor
}

// PGPool возвращает пул подключений к PostgreSQL (ленивая инициализация).
func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}
		if err = pool.Ping(ctx); err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}
	return d.pgPool
}

func (d *diContainer) TxManager(ctx context.Context) applicationPart.TxManager {
	if d.txManager == nil {
		m, err := trmmanager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать менеджер транзакций", "error", err)
			os.Exit(1)
		}
		d.txManager = m
	}
	return d.txManager
}

// PartRepository возвращает репозиторий деталей.
func (d *diContainer) PartRepository(ctx context.Context) applicationPart.PartRepository {
	if d.partRepo == nil {
		d.partRepo = partRepository.NewRepository(d.PGPool(ctx), trmpgx.DefaultCtxGetter)
	}
	return d.partRepo
}

// PartService возвращает сервис бизнес-логики деталей.
func (d *diContainer) PartService(ctx context.Context) inventoryAPIv1.PartService {
	if d.partSvc == nil {
		d.partSvc = applicationPart.NewService(d.PartRepository(ctx), domainService.NewCompatibilityChecker(), d.TxManager(ctx))
	}
	return d.partSvc
}

// InventoryAPI возвращает gRPC-обработчик InventoryService.
func (d *diContainer) InventoryAPI(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryHandler == nil {
		d.inventoryHandler = inventoryAPIv1.NewAPI(d.PartService(ctx))
	}
	return d.inventoryHandler
}

func (d *diContainer) IAmConn(_ context.Context) *grpc.ClientConn {
	if nil == d.iamCONN {
		conn, err := grpc.NewClient(

			config.AppConfig().IAMClient.Address,

			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			slog.Error("не удалось подключиться к iam сервису", "error", err)
			os.Exit(1)
		}
		closer.Add("grpc соединение с IamService", func(_ context.Context) error {
			return conn.Close()
		})

		d.iamCONN = conn
	}

	return d.iamCONN
}

func (d *diContainer) IAMClient(ctx context.Context) interceptor.IAMClient {
	if nil == d.iamClient {
		d.iamClient = iamGRPCClient.New(authv1.NewAuthServiceClient(d.IAmConn(ctx)))
	}

	return d.iamClient
}

func (d *diContainer) AuthInterceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	if nil == d.authInterceptor {
		d.authInterceptor = interceptor.NewAuth(d.IAMClient(ctx))
	}

	return d.authInterceptor
}
