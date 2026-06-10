package v1

import (
	"context"

	"github.com/1mpuser/inventory/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

type PartService interface {
	List(ctx context.Context, uuids []string, partType inventoryv1.PartType) ([]model.Part, error)
	Get(ctx context.Context, id string) (model.Part, error)
}
