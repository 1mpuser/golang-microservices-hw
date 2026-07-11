package assembly

import (
	"context"
	"math/rand"
	"time"

	"github.com/1mpuser/assembly/internal/model"
)

// Assemble обрабатывает оплату заказа: эмулирует сборку корабля и публикует ShipAssembled.
//
// TODO(неделя 5, часть 4): реализовать —
//  1. Сгенерировать случайное время сборки в диапазоне [min, max] и «поспать» его (эмуляция).
//  2. Собрать model.ShipAssembled (новый event_uuid, order_uuid и user_uuid из входного события,
//     build_time_sec, assembled_at = time.Now()).
//  3. Опубликовать через s.producer.ProduceShipAssembled.
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
