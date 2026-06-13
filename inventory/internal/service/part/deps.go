package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/record"
)

type PartRepository interface {
	ListPartsByUuids(ctx context.Context, uuids []uuid.UUID) ([]record.Part, error)
	ListPartsByPartType(ctx context.Context, partType model.PartType) ([]record.Part, error)
	ListAllParts(ctx context.Context) ([]record.Part, error)
	Get(ctx context.Context, uuid uuid.UUID) (record.Part, error)
}
