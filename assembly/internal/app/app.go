package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/1mpuser/assembly/internal/config"
	orderPaidConsumer "github.com/1mpuser/assembly/internal/consumer/order_paid"
	"github.com/1mpuser/platform/pkg/closer"
	"github.com/1mpuser/platform/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

// App управляет жизненным циклом AssemblyService (Kafka consumer + producer).
type App struct {
	di       *diContainer
	consumer *orderPaidConsumer.Consumer
}

// New создаёт приложение и инициализирует зависимости.
func New(ctx context.Context) *App {
	a := &App{}
	a.initDeps(ctx)
	return a
}

// Run запускает consumer OrderPaid, ждёт сигнала ОС и выполняет graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		slog.Info(
			"запуск AssemblyService: слушаем OrderPaid",
			"topic", config.AppConfig().OrderPaidConsumer.Topic,
			"group", config.AppConfig().OrderPaidConsumer.GroupID,
		)
		errCh <- a.consumer.Run(ctx)
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
		a.initConsumer,
	} {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.di = newDiContainer()
}

func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) initConsumer(ctx context.Context) {
	a.consumer = a.di.OrderPaidConsumer(ctx)
}
