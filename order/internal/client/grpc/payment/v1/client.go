package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/1mpuser/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/1mpuser/order/internal/errors"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

type client struct {
	grpcClient paymentv1.PaymentServiceClient
}

func New(grpcClient paymentv1.PaymentServiceClient) *client {
	return &client{
		grpcClient: grpcClient,
	}
}

func (c *client) PayOrder(ctx context.Context, orderId string, paymentMethod paymentv1.PaymentMethod) (*converter.TransactionUUID, error) {
	resp, err := c.grpcClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderId,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				return nil, errs.ErrInvalidPaymentMethod
			}
		}

		return nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	return converter.DTOToModel(resp.TransactionUuid), nil
}
