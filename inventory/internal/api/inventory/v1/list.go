package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/1mpuser/inventory/internal/api/convertor"
	errs "github.com/1mpuser/inventory/internal/errors"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func (a *api) List(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, req.Uuids, convertor.PartTypeFromProto(req.PartType))
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "детали не найдена с id: %s", req.GetUuids())
		}

		return nil, status.Errorf(codes.Internal, "ошибка получения деталей: %v", err)
	}

	dtos := make([]*inventoryv1.Part, 0, len(parts))

	for _, part := range parts {
		dtos = append(dtos, convertor.PartToDto(part))
	}

	return &inventoryv1.ListPartsResponse{
		Parts: dtos,
	}, nil
}
