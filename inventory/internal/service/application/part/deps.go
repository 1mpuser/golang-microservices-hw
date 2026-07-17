package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/record"
)

type CompatibilityChecker interface {
	Check(parts []model.Part) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type PartRepository interface {
	ListPartsByUuids(ctx context.Context, uuids []uuid.UUID) ([]record.PartRecord, error)
	ListPartsByPartType(ctx context.Context, partType model.PartType) ([]record.PartRecord, error)
	ListAllParts(ctx context.Context) ([]record.PartRecord, error)
	ListPartsForUpdate(ctx context.Context, uuids []uuid.UUID) ([]record.PartRecord, error)
	Get(ctx context.Context, uuid uuid.UUID) (record.PartRecord, error)
	UpdateParts(ctx context.Context, parts []model.Part) error
}
