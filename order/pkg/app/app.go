// Package app — тонкая обёртка над internal/app для API-тестов.
// Собирает зависимости OrderService так же, как internal/app/di.go,
// но принимает пул БД, менеджер транзакций и gRPC-клиентов напрямую.
package app

import (
	"net/http"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	orderAPIv1 "github.com/1mpuser/order/internal/api/order/v1"
	inventoryGRPCClient "github.com/1mpuser/order/internal/client/grpc/inventory/v1"
	paymentGRPCClient "github.com/1mpuser/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/1mpuser/order/internal/repository/order"
	orderService "github.com/1mpuser/order/internal/service/order"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler собирает зависимости OrderService и возвращает HTTP-обработчик.
func NewHTTPHandler(
	pool *pgxpool.Pool,
	txManager *trmmanager.Manager,
	invClient inventoryv1.InventoryServiceClient,
	payClient paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	orderRepo := orderRepository.NewRepository(pool, trmpgx.DefaultCtxGetter)

	svc := orderService.NewService(
		txManager,
		orderRepo,
		inventoryGRPCClient.New(invClient),
		paymentGRPCClient.New(payClient),
	)

	api := orderAPIv1.NewAPI(svc)

	return orderv1.NewServer(api)
}
