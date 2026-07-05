// Package app — тонкая обёртка над internal/app для API-тестов.
// Собирает зависимости InventoryService так же, как internal/app/di.go,
// но принимает пул БД напрямую (без чтения конфига).
package app

import (
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryAPIv1 "github.com/1mpuser/inventory/internal/api/inventory/v1"
	partRepository "github.com/1mpuser/inventory/internal/repository/part"
	applicationPart "github.com/1mpuser/inventory/internal/service/application/part"
	domainService "github.com/1mpuser/inventory/internal/service/domain"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

// Interceptors возвращает gRPC-опции сервера.
func Interceptors() []grpc.ServerOption {
	return nil
}

// RegisterServices собирает зависимости InventoryService и регистрирует их на gRPC-сервере.
func RegisterServices(server *grpc.Server, pool *pgxpool.Pool) {
	txManager, err := trmmanager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		panic(fmt.Sprintf("создать менеджер транзакций: %v", err))
	}

	repo := partRepository.NewRepository(pool, trmpgx.DefaultCtxGetter)
	svc := applicationPart.NewService(repo, domainService.NewCompatibilityChecker(), txManager)
	api := inventoryAPIv1.NewAPI(svc)

	inventoryv1.RegisterInventoryServiceServer(server, api)
}
