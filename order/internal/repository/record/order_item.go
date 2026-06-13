package record

import (
	"time"

	"github.com/1mpuser/order/internal/model"
	"github.com/google/uuid"
)

type OrderItem struct {
	OrderUUID uuid.UUID
	PartUUID  uuid.UUID
	PartType  model.PartType
	Price     int64
	CreatedAt time.Time
}
