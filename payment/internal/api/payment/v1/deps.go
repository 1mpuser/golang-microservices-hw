package v1

import (
	"context"

	"github.com/1mpuser/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, payRequest model.PayRequest) (string, error)
}
