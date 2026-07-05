package record

import (
	"time"

	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
)

type OrderItem struct {
	OrderUUID uuid.UUID
	PartUUID  uuid.UUID
	PartType  model.PartType
	Price     int64
	CreatedAt time.Time
}
