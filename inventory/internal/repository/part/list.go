package part

import (
	"context"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/convertor"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

func (r *repository) ListPartsByUuids(_ context.Context, uuids []uuid.UUID) ([]model.Part, error) {

	parts := make([]model.Part, 0, len(uuids))

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, id := range uuids {
		part, ok := r.data[id]

		if !ok {
			return nil, errs.ErrPartNotFound
		}

		parts = append(parts, convertor.PartToModel(part))

	}

	return parts, nil

}

func (r *repository) ListPartsByPartType(_ context.Context, partType inventoryv1.PartType) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0)

	for _, part := range r.data {
		if partType == inventoryv1.PartType_PART_TYPE_UNSPECIFIED || partType == part.PartType {
			parts = append(parts, convertor.PartToModel(part))
		}
	}

	return parts, nil

}
