package payment

import (
	"context"
	"log/slog"

	errs "github.com/1mpuser/payment/internal/errors"
	"github.com/1mpuser/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, payRequest model.PayRequest) (string, error) {
	_, err := uuid.Parse(payRequest.OrderUUID)

	if err != nil {
		return "", errs.ErrInvalidOrderUUID
	}

	if payRequest.PaymentMethod == model.PaymentMethodUnspecified {
		return "", errs.ErrInvalidPaymentMethod
	}

	transactionUuid := uuid.New()

	slog.Info(
		"оплата прошла успешно",
		"order_uuid", payRequest.OrderUUID,
		"transaction_uuid", transactionUuid,
	)

	return transactionUuid.String(), nil

}
