package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderAPI "github.com/1mpuser/order/internal/api/order/v1"
	inventoryClient "github.com/1mpuser/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/1mpuser/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/1mpuser/order/internal/repository/order"
	orderService "github.com/1mpuser/order/internal/service/order"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	httpAddress             = ":7884"

	grpcKeepaliveTime    = 10 * time.Second
	grpcKeepaliveTimeout = 3 * time.Second

	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 10 * time.Second
	httpWriteTimeout      = 10 * time.Second
	httpIdleTimeout       = 60 * time.Second
	httpShutdownTimeout   = 5 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		os.Exit(1)
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		if closeErr := inventoryConn.Close(); closeErr != nil {
			slog.Error("ошибка закрытия соединения с InventoryService", "error", closeErr)
		}
		os.Exit(1)
	}

	dbUri := os.Getenv("DB_URI")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, dbUri)

	if err != nil {
		slog.Error("не удалось подключиться к базе данных", "error", err)

		os.Exit(1)
	}

	defer pool.Close()

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))

	if err != nil {
		slog.Error("ошибка создания менеджера транзакций", "error", err)
		os.Exit(1)
	}

	txGetter := trmpgx.DefaultCtxGetter

	repo := orderRepository.NewRepository(pool, txGetter)

	invClient := inventoryClient.New(inventoryv1.NewInventoryServiceClient(inventoryConn))
	payClient := paymentClient.New(paymentv1.NewPaymentServiceClient(paymentConn))

	svc := orderService.NewService(txManager, repo, invClient, payClient)
	api := orderAPI.NewAPI(svc)

	orderServer, err := orderv1.NewServer(api)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		if closeErr := inventoryConn.Close(); closeErr != nil {
			slog.Error("ошибка закрытия соединения с InventoryService", "error", closeErr)
		}
		if closeErr := paymentConn.Close(); closeErr != nil {
			slog.Error("ошибка закрытия соединения с PaymentService", "error", closeErr)
		}
		os.Exit(1)
	}

	defer inventoryConn.Close()
	defer paymentConn.Close()

	httpServer := &http.Server{
		Addr:              httpAddress,
		Handler:           orderServer,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	go func() {
		slog.Info("запуск OrderService", "адрес", httpAddress)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ошибка запуска сервера", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("завершение работы OrderService...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("ошибка graceful shutdown", "error", err)
	}
}
