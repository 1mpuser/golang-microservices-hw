package part

import (
	"context"

	"github.com/1mpuser/inventory/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

type PartRepository interface {
	ListPartsByUuids(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error)
	ListPartsByPartType(ctx context.Context, partType inventoryv1.PartType) ([]model.Part, error)
	Get(ctx context.Context, uuid uuid.UUID) (model.Part, error)
}
