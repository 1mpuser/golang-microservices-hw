package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/order/internal/errors"
	orderService "github.com/1mpuser/order/internal/service/order"
	"github.com/1mpuser/order/internal/service/order/mocks"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
	)

	tests := []struct {
		name      string
		orderUUID string
		setupMock func(repo *mocks.OrderRepository)
		wantErr   error
	}{
		{
			name:      "успешная отмена заказа",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Delete(mock.Anything, orderUUID).
					Return(nil)
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
					Delete(mock.Anything, orderUUID).
					Return(errs.ErrOrderNotFound)
			},
			wantErr: errs.ErrOrderNotFound,
		},
		{
			name:      "заказ уже оплачен",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Delete(mock.Anything, orderUUID).
					Return(errs.ErrOrderAlreadyPaid)
			},
			wantErr: errs.ErrOrderAlreadyPaid,
		},
		{
			name:      "заказ уже отменён",
			orderUUID: orderUUID.String(),
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Delete(mock.Anything, orderUUID).
					Return(errs.ErrOrderCancelled)
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

			tc.setupMock(orderRepo)

			svc := orderService.NewService(txManager, orderRepo, inventoryClient, paymentClient)
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
