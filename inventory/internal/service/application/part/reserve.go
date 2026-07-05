package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/convertor"
)

func (s *service) Reserve(ctx context.Context, uuids []string) error {
	return s.updateReserved(ctx, uuids, func(p *model.Part) error {
		return p.Reserve(1)
	})
}

// updateReserved загружает детали по uuids в транзакции, применяет к каждой apply
// (Reserve/Release) и сохраняет обновлённое состояние одним батчем.
func (s *service) updateReserved(ctx context.Context, uuids []string, apply func(*model.Part) error) error {
	uuidsChecked := make([]uuid.UUID, 0, len(uuids))

	for _, id := range uuids {
		uid, err := uuid.Parse(id)
		if err != nil {
			return errs.ErrInvalidUUID
		}

		uuidsChecked = append(uuidsChecked, uid)
	}

	return s.txManager.Do(ctx, func(ctx context.Context) error {
		partRecords, err := s.partRepository.ListPartsByUuids(ctx, uuidsChecked)
		if err != nil {
			return err
		}

		if len(partRecords) != len(uuidsChecked) {
			return errs.ErrPartNotFound
		}

		parts := make([]model.Part, 0, len(partRecords))

		for _, partRecord := range partRecords {
			part, err := convertor.PartRecordToModel(partRecord)
			if err != nil {
				return err
			}

			if err := apply(&part); err != nil {
				return err
			}

			parts = append(parts, part)
		}

		return s.partRepository.UpdateParts(ctx, parts)
	})
}
