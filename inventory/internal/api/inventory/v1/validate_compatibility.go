package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

type slot struct {
	uuid     string
	partType model.PartType
}

func (a *api) ValidateCompatibility(ctx context.Context, req *inventoryv1.ValidateCompatibilityRequest) (*inventoryv1.ValidateCompatibilityResponse, error) {
	if req.GetHullUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "hull_uuid обязателен")
	}
	if req.GetEngineUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "engine_uuid обязателен")
	}

	// Обязательные и опциональные слоты с ожидаемым типом детали.
	slots := []slot{
		{req.GetHullUuid(), model.PartTypeHull},
		{req.GetEngineUuid(), model.PartTypeEngine},
	}
	if req.GetShieldUuid() != "" {
		slots = append(slots, slot{req.GetShieldUuid(), model.PartTypeShield})
	}
	if req.GetWeaponUuid() != "" {
		slots = append(slots, slot{req.GetWeaponUuid(), model.PartTypeWeapon})
	}

	seen := make(map[string]struct{}, len(slots))
	uuids := make([]string, 0, len(slots))

	// Проверяем каждый слот: нет дублей, деталь существует и её тип соответствует слоту.
	for _, s := range slots {
		if _, dup := seen[s.uuid]; dup {
			return nil, status.Errorf(codes.InvalidArgument, "дублирующийся uuid детали: %s", s.uuid)
		}
		seen[s.uuid] = struct{}{}
		uuids = append(uuids, s.uuid)

		part, err := a.partService.Get(ctx, s.uuid)
		if err != nil {
			switch {
			case errors.Is(err, errs.ErrInvalidUUID):
				return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", s.uuid)
			case errors.Is(err, errs.ErrPartNotFound):
				return nil, status.Errorf(codes.NotFound, "деталь не найдена: %s", s.uuid)
			default:
				return nil, status.Errorf(codes.Internal, "получить деталь: %v", err)
			}
		}

		if part.PartType() != s.partType {
			return nil, status.Errorf(codes.InvalidArgument, "деталь %s не подходит для слота %s (тип %s)", s.uuid, s.partType, part.PartType())
		}
	}

	if err := a.partService.ValidateCompatibility(ctx, uuids); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %v", err)
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "деталь не найдена: %v", err)
		case errors.Is(err, errs.ErrIncompatibleParts):
			return nil, status.Errorf(codes.FailedPrecondition, "детали несовместимы: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "ошибка проверки совместимости: %v", err)
		}
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
