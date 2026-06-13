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
	"github.com/1mpuser/order/internal/service/input"
	orderService "github.com/1mpuser/order/internal/service/order"
	"github.com/1mpuser/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		hullUUID   = uuid.New()
		engineUUID = uuid.New()

		errRepo = errors.New("ошибка хранилища")

		partsInStock = []model.Part{
			{UUID: hullUUID.String(), Name: "Hull", Price: 500000, PartType: model.PartTypeHull, StockQuantity: 10},
			{UUID: engineUUID.String(), Name: "Engine", Price: 300000, PartType: model.PartTypeEngine, StockQuantity: 5},
		}

		partsOutOfStock = []model.Part{
			{UUID: hullUUID.String(), Name: "Hull", Price: 500000, PartType: model.PartTypeHull, StockQuantity: 10},
			{UUID: engineUUID.String(), Name: "Engine", Price: 300000, PartType: model.PartTypeEngine, StockQuantity: 0},
		}
	)

	tests := []struct {
		name      string
		in        input.CreateOrderInput
		setupMock func(repo *mocks.OrderRepository, client *mocks.InventoryClient)
		wantErr   error
	}{
		{
			name: "успешное создание заказа",
			in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(mock.Anything, []string{hullUUID.String(), engineUUID.String()}).
					Return(partsInStock, nil)

				repo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(o record.Order) bool {
						return o.HullUUID == hullUUID &&
							o.EngineUUID == engineUUID &&
							o.TotalPrice == 800000 &&
							o.Status == model.OrderStatusPendingPayment
					}), mock.Anything).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "детали не найдены на складе",
			in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(_ *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(mock.Anything, []string{hullUUID.String(), engineUUID.String()}).
					Return(nil, nil)
			},
			wantErr: errs.ErrInventoryPartsNotFound,
		},
		{
			name: "ошибка получения деталей из inventory",
			in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(_ *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(mock.Anything, []string{hullUUID.String(), engineUUID.String()}).
					Return(nil, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
		{
			name: "деталь отсутствует на складе",
			in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(_ *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(mock.Anything, []string{hullUUID.String(), engineUUID.String()}).
					Return(partsOutOfStock, nil)
			},
			wantErr: errs.ErrOutOfStock,
		},
		{
			name: "ошибка при сохранении заказа",
			in: input.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(mock.Anything, []string{hullUUID.String(), engineUUID.String()}).
					Return(partsInStock, nil)

				repo.EXPECT().
					Create(mock.Anything, mock.Anything, mock.Anything).
					Return(errRepo)
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

			txManager.EXPECT().
				Do(mock.Anything, mock.Anything).
				RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				}).
				Maybe()

			tc.setupMock(orderRepo, inventoryClient)

			svc := orderService.NewService(txManager, orderRepo, inventoryClient, paymentClient)
			result, err := svc.Create(ctx, tc.in)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.OrderUUID)
				assert.Equal(t, int64(800000), result.TotalPrice)
			}
		})
	}
}
