package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/1mpuser/iam/internal/config"
	"github.com/1mpuser/iam/internal/interceptor"
	"github.com/1mpuser/platform/pkg/closer"
	"github.com/1mpuser/platform/pkg/grpc/health"
	"github.com/1mpuser/platform/pkg/logger"
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
	userv1 "github.com/1mpuser/shared/pkg/proto/user/v1"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 5 * time.Second
)

// App управляет жизненным циклом IAMService.
type App struct {
	di         *diContainer
	grpcServer *grpc.Server
	listener   net.Listener
}

// New создаёт и инициализирует приложение.
func New(ctx context.Context) *App {
	a := &App{}
	a.initDeps(ctx)
	return a
}

// Run запускает gRPC-сервер, ждёт сигнала ОС и выполняет graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("запуск IAMService", "адрес", config.AppConfig().GRPC.Address())
		errCh <- a.grpcServer.Serve(a.listener)
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
		a.initListener,
		a.initGRPCServer,
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

func (a *App) initListener(_ context.Context) {
	lis, err := net.Listen("tcp", config.AppConfig().GRPC.Address()) //nolint:noctx // адрес из конфига
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}
	a.listener = lis
}

func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = grpc.NewServer(
		// Единый error-interceptor: доменные ошибки errs.* → gRPC-коды.
		grpc.ChainUnaryInterceptor(interceptor.ErrorUnaryInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
	)

	// Дёргаем API-геттеры ДО регистрации closer'а сервера — это запускает
	// ленивую инициализацию (API → Service → Repository → PGPool/Redis),
	// которые регистрируют себя в closer раньше. Closer работает по LIFO:
	// сначала остановится gRPC, потом закроются пул БД и Redis.
	authAPI := a.di.AuthAPI(ctx)
	userAPI := a.di.UserAPI(ctx)

	closer.Add("gRPC сервер", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	authv1.RegisterAuthServiceServer(a.grpcServer, authAPI)
	userv1.RegisterUserServiceServer(a.grpcServer, userAPI)
}
