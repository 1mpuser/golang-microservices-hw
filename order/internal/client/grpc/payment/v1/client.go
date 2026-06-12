package v1

import (
	"context"

	"github.com/1mpuser/order/internal/client/grpc/payment/v1/converter"
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
		return nil, err
	}

	return converter.DTOToModel(resp.TransactionUuid), nil
}
