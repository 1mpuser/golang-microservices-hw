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

func (a *api) Get(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid обязателен")
	}

	part, err := a.partService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "деталь не найдена с id: %s", req.GetUuid())
		}
		return nil, status.Errorf(codes.Internal, "ошибка получения детали: %v", err)
	}

	return &inventoryv1.GetPartResponse{
		Part: convertor.PartToDto(part),
	}, nil
}
