package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/service/input"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
)

// CreateOrder создаёт новый заказ на постройку космического корабля.
func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	in := input.CreateOrderInput{
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
	}

	if shieldUUID, ok := req.GetShieldUUID().Get(); ok {
		in.ShieldUUID = &shieldUUID
	}

	if weaponUUID, ok := req.GetWeaponUUID().Get(); ok {
		in.WeaponUUID = &weaponUUID
	}

	result, err := a.orderService.Create(ctx, in)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInventoryPartsNotFound):
			return &orderv1.CreateOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrOutOfStock):
			return &orderv1.CreateOrderConflict{Code: http.StatusConflict, Message: err.Error()}, nil
		default:
			return nil, err
		}
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  result.OrderUUID,
		TotalPrice: result.TotalPrice,
	}, nil
}
