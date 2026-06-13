package record

import (
	"time"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
)

type Order struct {
	OrderUUID       uuid.UUID            `db:"uuid"`
	HullUUID        uuid.UUID            `db:"hull_uuid"`
	EngineUUID      uuid.UUID            `db:"engine_uuid"`
	ShieldUUID      *uuid.UUID           `db:"shield_uuid"` // опциональный
	WeaponUUID      *uuid.UUID           `db:"weapon_uuid"` // опциональный
	TotalPrice      int64                `db:"total_price"` // в копейках
	TransactionUUID *uuid.UUID           `db:"transaction_uuid"`
	PaymentMethod   *model.PaymentMethod `db:"payment_method"`
	Status          model.OrderStatus
	CreatedAt       time.Time `db:"created_at"`
}
