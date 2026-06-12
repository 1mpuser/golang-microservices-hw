package order

import (
	"context"
	"fmt"
	"time"

	"github.com/1mpuser/order/internal/converter"
	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	repositoryConvertor "github.com/1mpuser/order/internal/repository/converter"
	"github.com/1mpuser/order/internal/service/input"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, in input.CreateOrderInput) (*converter.CreateModelDto, error) {
	uuids := []string{in.HullUUID.String(), in.EngineUUID.String()}

	if in.ShieldUUID != nil {
		uuids = append(uuids, in.ShieldUUID.String())
	}
	if in.WeaponUUID != nil {
		uuids = append(uuids, in.WeaponUUID.String())
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	parts, err := s.inventoryClient.ListParts(ctx, uuids)

	if err != nil {
		return nil, fmt.Errorf("получить детали: %w", err)
	}

	if len(parts) == 0 {
		return nil, errs.ErrInventoryPartsNotFound
	}

	var totalPrice int64 = 0

	for _, part := range parts {
		if part.StockQuantity == 0 {
			return nil, errs.ErrOutOfStock
		}

		totalPrice += part.Price
	}

	orderUUID := uuid.New()

	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   in.HullUUID,
		EngineUUID: in.EngineUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	if in.ShieldUUID != nil {
		order.ShieldUUID = in.ShieldUUID
	}
	if in.WeaponUUID != nil {
		order.WeaponUUID = in.WeaponUUID
	}

	err = s.orderRepository.Create(ctx, repositoryConvertor.ModelToRecord(order))

	if err != nil {
		return nil, err
	}

	return &converter.CreateModelDto{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil

}
