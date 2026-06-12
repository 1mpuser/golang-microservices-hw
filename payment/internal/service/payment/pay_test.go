package payment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/payment/internal/errors"
	"github.com/1mpuser/payment/internal/model"
	paymentService "github.com/1mpuser/payment/internal/service/payment"
)

func TestPay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		req     model.PayRequest
		wantErr error
	}{
		{
			name: "успешная оплата картой",
			req: model.PayRequest{
				OrderUUID:     uuid.New().String(),
				PaymentMethod: model.PaymentMethodCard,
			},
			wantErr: nil,
		},
		{
			name: "успешная оплата через СБП",
			req: model.PayRequest{
				OrderUUID:     uuid.New().String(),
				PaymentMethod: model.PaymentMethodSBP,
			},
			wantErr: nil,
		},
		{
			name: "неверный формат order_uuid",
			req: model.PayRequest{
				OrderUUID:     "не-uuid",
				PaymentMethod: model.PaymentMethodCard,
			},
			wantErr: errs.ErrInvalidOrderUUID,
		},
		{
			name: "пустой order_uuid",
			req: model.PayRequest{
				OrderUUID:     "",
				PaymentMethod: model.PaymentMethodCard,
			},
			wantErr: errs.ErrInvalidOrderUUID,
		},
		{
			name: "метод оплаты не указан",
			req: model.PayRequest{
				OrderUUID:     uuid.New().String(),
				PaymentMethod: model.PaymentMethodUnspecified,
			},
			wantErr: errs.ErrInvalidPaymentMethod,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := paymentService.NewService()
			transactionUUID, err := svc.Pay(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, transactionUUID)
				return
			}

			require.NoError(t, err)
			_, err = uuid.Parse(transactionUUID)
			assert.NoError(t, err)
		})
	}
}
