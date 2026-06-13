package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	paymentConverter "github.com/1mpuser/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/repository/record"
	orderService "github.com/1mpuser/order/internal/service/order"
	"github.com/1mpuser/order/internal/service/order/mocks"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

func TestPay(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		orderUUID       = uuid.New()
		transactionUUID = uuid.New()

		errPaymentService = errors.New("ошибка сервиса оплаты")

		pendingOrder = record.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPendingPayment,
		}

		paidOrder = record.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPaid,
		}

		cancelledOrder = record.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusCancelled,
		}
	)

	tests := []struct {
		name      string
		orderUUID string
		setupMock func(repo *mocks.OrderRepository, client *mocks.PaymentClient)
		wantErr   error
	}{
		{
			name:      "успешная оплата",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&pendingOrder, nil)

				client.EXPECT().
					PayOrder(mock.Anything, orderUUID.String(), paymentv1.PaymentMethod_PAYMENT_METHOD_CARD).
					Return(paymentConverter.DTOToModel(transactionUUID.String()), nil)

				repo.EXPECT().
					Pay(mock.Anything, orderUUID, model.PaymentMethodCard, transactionUUID).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "заказ не найден",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(nil, errs.ErrOrderNotFound)
			},
			wantErr: errs.ErrOrderNotFound,
		},
		{
			name:      "заказ уже оплачен",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&paidOrder, nil)
			},
			wantErr: errs.ErrOrderAlreadyPaid,
		},
		{
			name:      "заказ отменён",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.PaymentClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&cancelledOrder, nil)
			},
			wantErr: errs.ErrOrderCancelled,
		},
		{
			name:      "ошибка сервиса оплаты",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&pendingOrder, nil)

				client.EXPECT().
					PayOrder(mock.Anything, orderUUID.String(), paymentv1.PaymentMethod_PAYMENT_METHOD_CARD).
					Return(nil, errPaymentService)
			},
			wantErr: errPaymentService,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(orderRepo, paymentClient)

			svc := orderService.NewService(orderRepo, inventoryClient, paymentClient)
			result, err := svc.Pay(ctx, tc.orderUUID, paymentv1.PaymentMethod_PAYMENT_METHOD_CARD)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, transactionUUID, result.TransactionUUID)
			}
		})
	}
}
