package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/1mpuser/payment/internal/converter"
	errs "github.com/1mpuser/payment/internal/errors"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid обязателен")
	}

	transationId, err := a.paymentService.Pay(ctx, converter.DtoToModel(req))
	if err != nil {
		if errors.Is(err, errs.ErrInvalidOrderUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "Платеж не прошел из-за неправильного uuid: %s", req.GetOrderUuid())
		}
		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Errorf(codes.InvalidArgument, "платеж не прошел из-за неверного метода: %s", req.GetPaymentMethod())
		}
		return nil, status.Errorf(codes.Internal, "ошибка прохода платежа")
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transationId,
	}, nil
}
