package record

import (
	"time"

	"github.com/google/uuid"

	"github.com/1mpuser/inventory/internal/model"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64          // в копейках
	PartType      model.PartType `db:"part_type"`
	StockQuantity int64          `db:"stock_quantity"`
	CreatedAt     time.Time      `db:"created_at"`
}
