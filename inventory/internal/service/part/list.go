package part

import (
	"cmp"
	"context"
	"slices"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

func (s *service) List(ctx context.Context, uuids []string, partType inventoryv1.PartType) ([]model.Part, error) {

	if len(uuids) > 0 {
		uuidsChecked := make([]uuid.UUID, 0, len(uuids))

		for _, id := range uuids {
			idValidated, err := uuid.Parse(id)

			if err != nil {
				return nil, errs.ErrInvalidFormat
			}

			uuidsChecked = append(uuidsChecked, idValidated)
		}

		return s.partRepository.ListPartsByUuids(ctx, uuidsChecked)
	}

	parts, err := s.partRepository.ListPartsByPartType(ctx, partType)

	if err != nil {
		return nil, err
	}

	slices.SortFunc(parts, func(a, b model.Part) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return parts, nil

}
