package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderAPIv1 "github.com/1mpuser/order/internal/api/order/v1"
	iamGRPCClient "github.com/1mpuser/order/internal/client/grpc/iam/v1"
	inventoryGRPCClient "github.com/1mpuser/order/internal/client/grpc/inventory/v1"
	paymentGRPCClient "github.com/1mpuser/order/internal/client/grpc/payment/v1"
	"github.com/1mpuser/order/internal/config"
	assemblyConsumer "github.com/1mpuser/order/internal/consumer/assembly_consumer"
	"github.com/1mpuser/order/internal/interceptor"
	"github.com/1mpuser/order/internal/middleware"
	orderProducer "github.com/1mpuser/order/internal/producer/order_producer"
	orderRepository "github.com/1mpuser/order/internal/repository/order"
	orderService "github.com/1mpuser/order/internal/service/order"
	"github.com/1mpuser/platform/pkg/closer"
	kafkaConsumer "github.com/1mpuser/platform/pkg/kafka/consumer"
	kafkaProducer "github.com/1mpuser/platform/pkg/kafka/producer"
	kafkamw "github.com/1mpuser/platform/pkg/middleware/kafka"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
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

	iamConn   *grpc.ClientConn
	iamClient middleware.IAMClient

	shipConsumerGroup sarama.ConsumerGroup
	shipConsumer      *assemblyConsumer.Service

	syncProducer  sarama.SyncProducer
	orderProducer orderService.OrderProducer
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
			grpc.WithChainUnaryInterceptor(interceptor.SessionForwarder),
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

func (d *diContainer) IAMClient(ctx context.Context) middleware.IAMClient {
	if d.iamClient == nil {
		d.iamClient = iamGRPCClient.New(authv1.NewAuthServiceClient(d.IAMConn(ctx)))
	}
	return d.iamClient
}

func (d *diContainer) PaymentClient(ctx context.Context) orderService.PaymentClient {
	if d.payClient == nil {
		d.payClient = paymentGRPCClient.New(paymentv1.NewPaymentServiceClient(d.PaymentConn(ctx)))
	}
	return d.payClient
}

// SyncProducer — sarama SyncProducer для публикации событий OrderPaid.
func (d *diContainer) SyncProducer(_ context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		// Обязательно для SyncProducer — иначе SendMessage зависнет навсегда.
		cfg.Producer.Return.Successes = true
		cfg.Producer.RequiredAcks = sarama.WaitForAll

		sp, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers, cfg)
		if err != nil {
			slog.Error("не удалось создать Kafka sync producer", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka sync producer", func(_ context.Context) error {
			return sp.Close()
		})

		d.syncProducer = sp
	}
	return d.syncProducer
}

// OrderPaidProducer — доменный продюсер события OrderPaid.
func (d *diContainer) OrderPaidProducer(ctx context.Context) orderService.OrderProducer {
	if d.orderProducer == nil {
		platformProducer := kafkaProducer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig().OrderPaidProducer.Topic,
		)
		d.orderProducer = orderProducer.New(platformProducer)
	}
	return d.orderProducer
}

func (d *diContainer) OrderService(ctx context.Context) orderAPIv1.OrderService {
	if d.orderSvc == nil {
		d.orderSvc = orderService.NewService(
			d.TxManager(ctx),
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderPaidProducer(ctx),
		)
	}
	return d.orderSvc
}

// ShipAssembledConsumerGroup — sarama ConsumerGroup для чтения топика ShipAssembled.
func (d *diContainer) ShipAssembledConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if d.shipConsumerGroup == nil {
		cfg := sarama.NewConfig()
		// Читать с начала, если у группы ещё нет закоммиченного оффсета.
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().ShipAssembledConsumer.GroupID,
			cfg,
		)
		if err != nil {
			slog.Error("не удалось создать Kafka consumer group (ShipAssembled)", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka consumer group (ShipAssembled)", func(_ context.Context) error {
			return group.Close()
		})

		d.shipConsumerGroup = group
	}
	return d.shipConsumerGroup
}

func (d *diContainer) IAMConn(_ context.Context) *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time: grpcKeepaliveTime, Timeout: grpcKeepaliveTimeout, PermitWithoutStream: true,
			}),
		)
		if err != nil {
			slog.Error("не удалось подключиться к IAMService", "error", err)
			os.Exit(1)
		}
		closer.Add("gRPC соединение с IAMService", func(_ context.Context) error {
			return conn.Close()
		})
		d.iamConn = conn
	}
	return d.iamConn
}

// ShipAssembledConsumer — consumer события ShipAssembled (переводит заказ в ASSEMBLED).
func (d *diContainer) ShipAssembledConsumer(ctx context.Context) *assemblyConsumer.Service {
	if d.shipConsumer == nil {
		platformConsumer := kafkaConsumer.NewConsumer(
			d.ShipAssembledConsumerGroup(ctx),
			[]string{config.AppConfig().ShipAssembledConsumer.Topic},
			kafkaConsumer.WithMiddlewares(kafkamw.ConsumerLogging()),
		)

		d.shipConsumer = assemblyConsumer.NewService(
			platformConsumer,
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.TxManager(ctx),
		)
	}
	return d.shipConsumer
}

func (d *diContainer) HTTPHandler(ctx context.Context) http.Handler {
	if d.httpHandler == nil {
		api := orderAPIv1.NewAPI(d.OrderService(ctx))
		handler, err := orderv1.NewServer(api)
		if err != nil {
			slog.Error("не удалось создать HTTP-сервер", "error", err)
			os.Exit(1)
		}
		d.httpHandler = middleware.NewAuth(d.IAMClient(ctx))(handler)
	}
	return d.httpHandler
}
