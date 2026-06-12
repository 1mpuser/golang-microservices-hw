package record

import (
	"time"

	"github.com/1mpuser/order/internal/model"
	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *model.PaymentMethod
	Status          model.OrderStatus
	CreatedAt       time.Time
}
