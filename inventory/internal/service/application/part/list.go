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
		return s.listByUUIDs(ctx, uuids)
	}

	return s.listByPartType(ctx, partType)
}

// listByUUIDs возвращает детали по списку UUID. Все запрошенные (уникальные)
// детали должны существовать, иначе — ErrPartNotFound.
func (s *service) listByUUIDs(ctx context.Context, uuids []string) ([]model.Part, error) {
	uuidsChecked := make([]uuid.UUID, 0, len(uuids))
	uniqueUUIDs := make(map[uuid.UUID]struct{}, len(uuids))

	for _, id := range uuids {
		idValidated, err := uuid.Parse(id)
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}

		uuidsChecked = append(uuidsChecked, idValidated)
		uniqueUUIDs[idValidated] = struct{}{}
	}

	records, err := s.partRepository.ListPartsByUuids(ctx, uuidsChecked)
	if err != nil {
		return nil, err
	}

	if len(records) != len(uniqueUUIDs) {
		return nil, errs.ErrPartNotFound
	}

	return recordsToModels(records)
}

// listByPartType возвращает все детали или отфильтрованные по типу, отсортированные по имени.
func (s *service) listByPartType(ctx context.Context, partType model.PartType) ([]model.Part, error) {
	var (
		records []record.PartRecord
		err     error
	)

	if partType == model.PartTypeUnspecified {
		records, err = s.partRepository.ListAllParts(ctx)
	} else {
		records, err = s.partRepository.ListPartsByPartType(ctx, partType)
	}

	if err != nil {
		return nil, err
	}

	parts, err := recordsToModels(records)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(parts, func(a, b model.Part) int {
		return cmp.Compare(a.Name(), b.Name())
	})

	return parts, nil
}

func recordsToModels(records []record.PartRecord) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(records))

	for _, rec := range records {
		part, err := convertor.PartRecordToModel(rec)
		if err != nil {
			return nil, errs.ErrIncompatibleParts
		}

		parts = append(parts, part)
	}

	return parts, nil
}
