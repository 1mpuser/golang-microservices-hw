package order_test

import (
	"context"
	"errors"
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

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		orderUUID = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		orderRecord = record.Order{
			OrderUUID:  orderUUID,
			HullUUID:   uuid.New(),
			EngineUUID: uuid.New(),
			TotalPrice: 800000,
			Status:     model.OrderStatusPendingPayment,
		}
	)

	tests := []struct {
		name      string
		orderUUID string
		setupMock func(repo *mocks.OrderRepository)
		wantErr   error
	}{
		{
			name:      "заказ найден",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(&orderRecord, nil)
			},
			wantErr: nil,
		},
		{
			name:      "неверный формат uuid",
			orderUUID: "не-uuid",
			setupMock: func(_ *mocks.OrderRepository) {},
			wantErr:   errs.ErrInvalidUUID,
		},
		{
			name:      "заказ не найден",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(nil, errs.ErrOrderNotFound)
			},
			wantErr: errs.ErrOrderNotFound,
		},
		{
			name:      "ошибка репозитория",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)
			txManager := mocks.NewTxManager(t)

			tc.setupMock(orderRepo)

			svc := orderService.NewService(txManager, orderRepo, inventoryClient, paymentClient)
			order, err := svc.Get(ctx, tc.orderUUID)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, uuid.Nil, order.OrderUUID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, orderUUID, order.OrderUUID)
				assert.Equal(t, int64(800000), order.TotalPrice)
			}
		})
	}
}
