package v1

import (
	"context"

	"github.com/1mpuser/order/internal/converter"
	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/service/input"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

type OrderService interface {
	Create(ctx context.Context, in input.CreateOrderInput) (*converter.CreateModelDto, error)
	Cancel(ctx context.Context, orderUuid string) error
	Get(ctx context.Context, orderId string) (model.Order, error)
	Pay(ctx context.Context, orderUuid string, paymentMethod paymentv1.PaymentMethod) (*converter.PayDto, error)
}
