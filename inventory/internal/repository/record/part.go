package record

import (
	"time"

	"github.com/google/uuid"

	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64                // в копейках
	PartType      inventoryv1.PartType `db:"part_type"`
	StockQuantity int64                `db:"stock_quantity"`
	CreatedAt     time.Time            `db:"created_at"`
}
