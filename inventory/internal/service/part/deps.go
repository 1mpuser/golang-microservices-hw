package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/1mpuser/inventory/internal/repository/record"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

type PartRepository interface {
	ListPartsByUuids(ctx context.Context, uuids []uuid.UUID) ([]record.Part, error)
	ListPartsByPartType(ctx context.Context, partType inventoryv1.PartType) ([]record.Part, error)
	ListAllParts(ctx context.Context) ([]record.Part, error)
	Get(ctx context.Context, uuid uuid.UUID) (record.Part, error)
}
