package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/convertor"
)

func (s *service) ValidateCompatibility(ctx context.Context, uuids []string) error {
	if len(uuids) == 0 {
		return nil
	}

	uuidsChecked := make([]uuid.UUID, 0, len(uuids))

	for _, id := range uuids {
		idValidated, err := uuid.Parse(id)
		if err != nil {
			return errs.ErrInvalidUUID
		}

		uuidsChecked = append(uuidsChecked, idValidated)
	}

	records, err := s.partRepository.ListPartsByUuids(ctx, uuidsChecked)
	if err != nil {
		return fmt.Errorf("получить детали: %w", err)
	}

	parts := make([]model.Part, 0, len(records))

	for _, record := range records {
		part, err := convertor.PartRecordToModel(record)
		if err != nil {
			return errs.ErrIncompatibleParts
		}

		parts = append(parts, part)
	}

	return s.checker.Check(parts)
}
