package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/order/internal/errors"
	"github.com/1mpuser/order/internal/repository/record"
)

func (r *repository) Get(_ context.Context, id uuid.UUID) (record.Order, error) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	order, ok := r.data[id]

	if !ok {
		return record.Order{}, errs.ErrOrderNotFound
	}

	return order, nil
}
