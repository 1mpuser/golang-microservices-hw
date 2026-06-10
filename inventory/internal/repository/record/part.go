package record

import (
	"time"

	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     time.Time
}
