package model

import (
	"time"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      PartType
	StockQuantity int64
	CreatedAt     *time.Time
}
