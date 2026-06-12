package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/1mpuser/order/internal/converter"
	errs "github.com/1mpuser/order/internal/errors"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
)

// PayOrder проводит оплату заказа.
func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	paymentMethod := converter.PaymentMethodFromOpenAPI(req.GetPaymentMethod())

	result, err := a.orderService.Pay(ctx, params.OrderUUID.String(), paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.PayOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrOrderAlreadyPaid), errors.Is(err, errs.ErrOrderCancelled):
			return &orderv1.PayOrderConflict{Code: http.StatusConflict, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrInvalidPaymentMethod):
			return &orderv1.PayOrderBadRequest{Code: http.StatusBadRequest, Message: err.Error()}, nil
		default:
			return nil, err
		}
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: result.TransactionUUID,
	}, nil
}
