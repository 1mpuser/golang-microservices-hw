package record

import (
	"time"
)

type PartRecord struct {
	UUID          string
	Name          string
	Description   string
	Price         int64  // в копейках
	PartType      string `db:"part_type"`
	StockQuantity int    `db:"stock_quantity"`
	Reserved      int
	Properties    []byte
	CreatedAt     time.Time `db:"created_at"`
}
