package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/1mpuser/order/internal/converter"
	errs "github.com/1mpuser/order/internal/errors"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
)

// GetOrder возвращает полную информацию о заказе на постройку космического корабля.
func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return &orderv1.GetOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
		}

		return nil, err
	}

	return converter.OrderToDTO(order), nil
}
