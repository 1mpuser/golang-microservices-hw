package v1

import (
	"context"

	"github.com/1mpuser/inventory/internal/model"
)

type PartService interface {
	List(ctx context.Context, uuids []string, partType model.PartType) ([]model.Part, error)
	Get(ctx context.Context, id string) (model.Part, error)
}
