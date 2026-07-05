package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderAPIv1 "github.com/1mpuser/order/internal/api/order/v1"
	inventoryGRPCClient "github.com/1mpuser/order/internal/client/grpc/inventory/v1"
	paymentGRPCClient "github.com/1mpuser/order/internal/client/grpc/payment/v1"
	"github.com/1mpuser/order/internal/config"
	orderRepository "github.com/1mpuser/order/internal/repository/order"
	orderService "github.com/1mpuser/order/internal/service/order"
	"github.com/1mpuser/platform/pkg/closer"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

const (
	grpcKeepaliveTime    = 10 * time.Second
	grpcKeepaliveTimeout = 3 * time.Second
)

type diContainer struct {
	pgPool        *pgxpool.Pool
	txManager     orderService.TxManager
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
	orderRepo     orderService.OrderRepository
	invClient     orderService.InventoryClient
	payClient     orderService.PaymentClient
	orderSvc      orderAPIv1.OrderService
	httpHandler   http.Handler
}

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

func (d *diContainer) TxManager(ctx context.Context) orderService.TxManager {
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

func (d *diContainer) InventoryConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time: grpcKeepaliveTime, Timeout: grpcKeepaliveTimeout, PermitWithoutStream: true,
			}),
		)
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			os.Exit(1)
		}
		closer.Add("gRPC соединение с InventoryService", func(_ context.Context) error {
			return conn.Close()
		})
		d.inventoryConn = conn
	}
	return d.inventoryConn
}

func (d *diContainer) PaymentConn(_ context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time: grpcKeepaliveTime, Timeout: grpcKeepaliveTimeout, PermitWithoutStream: true,
			}),
		)
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			os.Exit(1)
		}
		closer.Add("gRPC соединение с PaymentService", func(_ context.Context) error {
			return conn.Close()
		})
		d.paymentConn = conn
	}
	return d.paymentConn
}

func (d *diContainer) OrderRepository(ctx context.Context) orderService.OrderRepository {
	if d.orderRepo == nil {
		d.orderRepo = orderRepository.NewRepository(d.PGPool(ctx), trmpgx.DefaultCtxGetter)
	}
	return d.orderRepo
}

func (d *diContainer) InventoryClient(ctx context.Context) orderService.InventoryClient {
	if d.invClient == nil {
		d.invClient = inventoryGRPCClient.New(inventoryv1.NewInventoryServiceClient(d.InventoryConn(ctx)))
	}
	return d.invClient
}

func (d *diContainer) PaymentClient(ctx context.Context) orderService.PaymentClient {
	if d.payClient == nil {
		d.payClient = paymentGRPCClient.New(paymentv1.NewPaymentServiceClient(d.PaymentConn(ctx)))
	}
	return d.payClient
}

func (d *diContainer) OrderService(ctx context.Context) orderAPIv1.OrderService {
	if d.orderSvc == nil {
		d.orderSvc = orderService.NewService(
			d.TxManager(ctx),
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}
	return d.orderSvc
}

func (d *diContainer) HTTPHandler(ctx context.Context) http.Handler {
	if d.httpHandler == nil {
		api := orderAPIv1.NewAPI(d.OrderService(ctx))
		handler, err := orderv1.NewServer(api)
		if err != nil {
			slog.Error("не удалось создать HTTP-сервер", "error", err)
			os.Exit(1)
		}
		d.httpHandler = handler
	}
	return d.httpHandler
}
