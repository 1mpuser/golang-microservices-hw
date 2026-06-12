package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/repository/record"
)

type repository struct {
	data map[uuid.UUID]record.Order
	mu   sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		data: make(map[uuid.UUID]record.Order),
	}
}
