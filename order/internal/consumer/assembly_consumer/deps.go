package assemblyconsumer

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/repository/record"
)

// OrderRepository — доступ к заказам (с блокировкой строки и сменой статуса).
type OrderRepository interface {
	GetForUpdate(ctx context.Context, id uuid.UUID) (*record.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
}

// InventoryClient — списание деталей со склада после сборки.
type InventoryClient interface {
	CommitParts(ctx context.Context, uuids []string) error
}

// TxManager — менеджер транзакций.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
