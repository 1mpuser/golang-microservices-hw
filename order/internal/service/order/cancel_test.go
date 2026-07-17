package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	"github.com/1mpuser/order/internal/repository/record"
	orderService "github.com/1mpuser/order/internal/service/order"
	"github.com/1mpuser/order/internal/service/order/mocks"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		hullUUID  = uuid.New()
		engUUID   = uuid.New()

		pendingOrder = record.Order{
			OrderUUID:  orderUUID,
			HullUUID:   hullUUID,
			EngineUUID: engUUID,
			Status:     model.OrderStatusPendingPayment,
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
		setupMock func(repo *mocks.OrderRepository, client *mocks.InventoryClient)
		wantErr   error
	}{
		{
			name:      "успешная отмена заказа",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&pendingOrder, nil)

				client.EXPECT().
					ReleaseParts(mock.Anything, mock.Anything).
					Return(nil)

				repo.EXPECT().
					UpdateStatus(mock.Anything, orderUUID, model.OrderStatusCancelled).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "неверный формат uuid",
			orderUUID: "не-uuid",
			setupMock: func(_ *mocks.OrderRepository, _ *mocks.InventoryClient) {},
			wantErr:   errs.ErrInvalidUUID,
		},
		{
			name:      "заказ не найден",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.InventoryClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(nil, errs.ErrOrderNotFound)
			},
			wantErr: errs.ErrOrderNotFound,
		},
		{
			name:      "заказ уже оплачен",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.InventoryClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&paidOrder, nil)
			},
			wantErr: errs.ErrOrderAlreadyPaid,
		},
		{
			name:      "заказ уже отменён",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository, _ *mocks.InventoryClient) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&cancelledOrder, nil)
			},
			wantErr: errs.ErrOrderCancelled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)
			txManager := mocks.NewTxManager(t)

			txManager.EXPECT().
				Do(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				}).
				Maybe()

			tc.setupMock(orderRepo, inventoryClient)

			svc := orderService.NewService(txManager, orderRepo, inventoryClient, paymentClient, mocks.NewOrderProducer(t))
			err := svc.Cancel(ctx, tc.orderUUID)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
