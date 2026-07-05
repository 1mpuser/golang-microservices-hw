package part

import (
	"context"

	"github.com/1mpuser/inventory/internal/model"
)

func (s *service) Release(ctx context.Context, uuids []string) error {
	return s.updateReserved(ctx, uuids, func(p *model.Part) error {
		return p.Release(1)
	})
}
