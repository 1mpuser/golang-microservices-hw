package v1

import (
	"context"
	"fmt"

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
		grpcClient,
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

func (c *client) ValidateCompatibility(ctx context.Context, hullUUID, engineUUID string, shieldUUID, weaponUUID *string) error {
	req := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   hullUUID,
		EngineUuid: engineUUID,
	}
	if shieldUUID != nil {
		req.ShieldUuid = *shieldUUID
	}
	if weaponUUID != nil {
		req.WeaponUuid = *weaponUUID
	}

	_, err := c.grpcClient.ValidateCompatibility(ctx, req)

	return mapInventoryError(err)
}

func (c *client) ReserveParts(ctx context.Context, uuids []string) error {
	_, err := c.grpcClient.ReserveParts(ctx, &inventoryv1.ReservePartsRequest{
		Uuids: uuids,
	})

	return mapInventoryError(err)
}

func (c *client) ReleaseParts(ctx context.Context, uuids []string) error {
	_, err := c.grpcClient.ReleaseParts(ctx, &inventoryv1.ReleasePartsRequest{
		Uuids: uuids,
	})

	return mapInventoryError(err)
}

// mapInventoryError переводит gRPC-ошибки InventoryService в доменные ошибки OrderService.
func mapInventoryError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("склад: %w", err)
	}

	switch st.Code() {
	case codes.NotFound:
		return errs.ErrPartNotFound
	case codes.FailedPrecondition:
		return errs.ErrIncompatibleParts
	case codes.ResourceExhausted:
		return errs.ErrOutOfStock
	case codes.InvalidArgument:
		return errs.ErrInvalidUUID
	default:
		return fmt.Errorf("склад: %w", err)
	}
}
