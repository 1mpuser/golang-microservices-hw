package order

import (
	"context"

	"github.com/1mpuser/order/internal/repository/record"
)

func (r *repository) Create(_ context.Context, order record.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[order.OrderUUID] = order

	return nil
}
