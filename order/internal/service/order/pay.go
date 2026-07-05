package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/converter"
	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

func (s *service) Pay(ctx context.Context, orderUuid string, paymentMethod paymentv1.PaymentMethod) (*converter.PayDto, error) {
	orderValidUuid, err := uuid.Parse(orderUuid)
	if err != nil {
		return nil, errs.ErrOrderNotFound
	}

	var result *converter.PayDto

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.orderRepository.Get(ctx, orderValidUuid)
		if err != nil {
			return err
		}

		switch order.Status {
		case model.OrderStatusPaid:
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			return errs.ErrOrderCancelled
		}

		transaction, err := s.paymentClient.PayOrder(ctx, orderUuid, paymentMethod)
		if err != nil {
			return err
		}

		transactionUUID, err := uuid.Parse(transaction.TransactionUUID)
		if err != nil {
			return err
		}

		err = s.orderRepository.Pay(ctx, order.OrderUUID, converter.PaymentMethodToModel(paymentMethod), transactionUUID)
		if err != nil {
			return err
		}

		result = converter.PayModelToDto(transactionUUID)

		return nil
	})

	return result, err
}
