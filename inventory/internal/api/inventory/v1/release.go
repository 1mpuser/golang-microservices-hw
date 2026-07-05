package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/1mpuser/inventory/internal/errors"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func (a *api) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	if len(req.GetUuids()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "uuids обязателен")
	}

	if err := a.partService.Release(ctx, req.GetUuids()); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %v", err)
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "деталь не найдена: %v", err)
		case errors.Is(err, errs.ErrNothingToRelease):
			return nil, status.Errorf(codes.FailedPrecondition, "нечего освобождать: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "ошибка освобождения деталей: %v", err)
		}
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
