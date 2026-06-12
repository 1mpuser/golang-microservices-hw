package order

import (
	"sync"

	"github.com/1mpuser/order/internal/repository/record"
	"github.com/google/uuid"
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
