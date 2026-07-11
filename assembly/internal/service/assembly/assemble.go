package assembly

import (
	"context"
	"math/rand"
	"time"

	"github.com/1mpuser/assembly/internal/model"
)

// Assemble обрабатывает оплату заказа: эмулирует сборку корабля и публикует ShipAssembled.
func (s *service) Assemble(ctx context.Context, event model.OrderPaid) error {
	amountOfBuildingSeconds := rand.Int63n(11) + 5

	sleepTime := time.Duration(amountOfBuildingSeconds) * time.Second

	time.Sleep(sleepTime)

	resultModel := model.ShipAssembled{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: amountOfBuildingSeconds,
		AssembledAt:  time.Now(),
	}

	return s.producer.ProduceShipAssembled(ctx, resultModel)
}
