package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/1mpuser/order/internal/config"
	assemblyConsumer "github.com/1mpuser/order/internal/consumer/assembly_consumer"
	"github.com/1mpuser/platform/pkg/closer"
	"github.com/1mpuser/platform/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

// App управляет жизненным циклом OrderService.
type App struct {
	di         *diContainer
	httpServer *http.Server
	consumer   *assemblyConsumer.Service
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context) *App {
	a := &App{}
	a.initDeps(ctx)
	return a
}

// Run запускает HTTP-сервер, ждёт сигнала ОС и выполняет graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 2)
	go func() {
		slog.Info("запуск OrderService", "адрес", config.AppConfig().HTTP.Address())
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	go func() {
		slog.Info(
			"запуск consumer ShipAssembled",
			"topic", config.AppConfig().ShipAssembledConsumer.Topic,
			"group", config.AppConfig().ShipAssembledConsumer.GroupID,
		)
		if err := a.consumer.RunConsumer(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
	}()

	var runErr error
	select {
	case runErr = <-errCh:
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка при завершении работы", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

func (a *App) initDeps(ctx context.Context) {
	for _, f := range []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initHTTPServer,
		a.initConsumer,
	} {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.di = &diContainer{}
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) initConsumer(ctx context.Context) {
	a.consumer = a.di.ShipAssembledConsumer(ctx)
}

func (a *App) initHTTPServer(ctx context.Context) {
	// Вызов HTTPHandler запускает цепочку ленивой инициализации:
	// HTTPHandler → OrderAPI → OrderService → OrderRepo + gRPC-клиенты → PGPool
	// PGPool и gRPC-соединения регистрируются в closer первыми.
	// HTTP-сервер регистрируется последним — при shutdown закроется первым (LIFO).
	handler := a.di.HTTPHandler(ctx)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	closer.Add("HTTP сервер", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
}
