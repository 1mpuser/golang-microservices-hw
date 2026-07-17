package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/client/grpc/payment/v1/converter"
	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/repository/record"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrderRepository interface {
	Create(ctx context.Context, order record.Order, orderItems []record.OrderItem) error
	Pay(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*record.Order, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*record.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
	ValidateCompatibility(ctx context.Context, hullUUID, engineUUID string, shieldUUID, weaponUUID *string) error
	ReserveParts(ctx context.Context, uuids []string) error
	ReleaseParts(ctx context.Context, uuids []string) error
	CommitParts(ctx context.Context, uuids []string) error
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderId string, paymentMethod paymentv1.PaymentMethod) (*converter.TransactionUUID, error)
}

type OrderProducer interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}
