package assembly

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/1mpuser/assembly/internal/model"
)

// Assemble обрабатывает оплату заказа: эмулирует сборку корабля и публикует ShipAssembled.
func (s *service) Assemble(ctx context.Context, event model.OrderPaid) error {
	buildSeconds := s.buildDurationSeconds()

	if buildSeconds > 0 {
		timer := time.NewTimer(time.Duration(buildSeconds) * time.Second)
		defer timer.Stop()

		// Ждём «сборку», но прерываемся на отмену контекста (graceful shutdown).
		select {
		case <-timer.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	resultModel := model.ShipAssembled{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildSeconds,
		AssembledAt:  time.Now(),
	}

	return s.producer.ProduceShipAssembled(ctx, resultModel)
}

// buildDurationSeconds возвращает эмулируемое время сборки в секундах в диапазоне
// [min, max]. При равных границах (в т.ч. 0/0) возвращает min без обращения к rand.
func (s *service) buildDurationSeconds() int64 {
	if s.maxBuildSeconds <= s.minBuildSeconds {
		return s.minBuildSeconds
	}

	//nolint:gosec // псевдослучайность для эмуляции времени сборки — не криптография
	return s.minBuildSeconds + rand.Int64N(s.maxBuildSeconds-s.minBuildSeconds+1)
}
