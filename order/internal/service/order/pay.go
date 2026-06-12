package order

import (
	"context"

	"github.com/1mpuser/order/internal/converter"
	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, orderUuid string, paymentMethod paymentv1.PaymentMethod) (*converter.PayDto, error) {

	orderValidUuid, err := uuid.Parse(orderUuid)

	if err != nil {
		return nil, errs.ErrOrderNotFound
	}

	order, err := s.orderRepository.Get(ctx, orderValidUuid)

	if err != nil {
		return nil, errs.ErrOrderNotFound
	}

	if order.Status != model.OrderStatusPendingPayment {
		return nil, errs.ErrOrderAlreadyPaid
	}

	transaction, err := s.paymentClient.PayOrder(ctx, orderUuid, paymentMethod)

	if err != nil {
		return nil, err
	}

	transactionUUID, err := uuid.Parse(transaction.TransactionUUID)

	if err != nil {
		return nil, err
	}

	err = s.orderRepository.Pay(ctx, order.OrderUUID, converter.PaymentMethodToModel(paymentMethod), transactionUUID)

	if err != nil {
		return nil, err
	}

	return converter.PayModelToDto(transactionUUID), nil
}
