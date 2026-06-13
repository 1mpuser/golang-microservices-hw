package part

import (
	"cmp"
	"context"
	"slices"

	"github.com/google/uuid"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/convertor"
	"github.com/1mpuser/inventory/internal/repository/record"
)

func (s *service) List(ctx context.Context, uuids []string, partType model.PartType) ([]model.Part, error) {
	if len(uuids) > 0 {
		uuidsChecked := make([]uuid.UUID, 0, len(uuids))

		for _, id := range uuids {
			idValidated, err := uuid.Parse(id)
			if err != nil {
				return nil, errs.ErrInvalidFormat
			}

			uuidsChecked = append(uuidsChecked, idValidated)
		}

		records, err := s.partRepository.ListPartsByUuids(ctx, uuidsChecked)

		if err != nil {
			return nil, err
		}

		parts := make([]model.Part, 0, len(records))

		for _, record := range records {
			parts = append(parts, convertor.PartToModel(record))
		}

		return parts, nil

	}

	var parts []record.Part
	var err error

	if partType == model.PartTypeUnspecified {
		parts, err = s.partRepository.ListAllParts(ctx)

		if err != nil {
			return nil, err
		}
	} else {
		parts, err = s.partRepository.ListPartsByPartType(ctx, partType)

		if err != nil {
			return nil, err
		}
	}

	partModels := make([]model.Part, 0, len(parts))

	for _, partRecord := range parts {
		partModels = append(partModels, convertor.PartToModel(partRecord))
	}

	slices.SortFunc(partModels, func(a, b model.Part) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return partModels, nil
}
