package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/1mpuser/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

type client struct {
	grpcClient inventoryv1.InventoryServiceClient
}

func New(grpcClient inventoryv1.InventoryServiceClient) *client {
	return &client{
		grpcClient: grpcClient,
	}
}

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.grpcClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids:    uuids,
		PartType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, errs.ErrInventoryPartsNotFound
			}
		}

		return nil, errs.ErrInventoryUnavailable
	}

	return converter.DTOToModel(resp), nil
}
