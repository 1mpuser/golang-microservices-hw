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

	orderHandler "github.com/1mpuser/order/pkg/handler"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	httpAddress             = ":8080"

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
	defer inventoryConn.Close()

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
		os.Exit(1)
	}
	defer paymentConn.Close()

	store := orderHandler.NewOrderStore()
	h := orderHandler.NewHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		store,
	)

	orderServer, err := orderHandler.SetupServer(h)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:              httpAddress,
		Handler:           orderServer,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
