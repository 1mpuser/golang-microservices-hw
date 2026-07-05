package domain

import (
	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
)

type compatibilityChecker struct{}

func NewCompatibilityChecker() *compatibilityChecker {
	return &compatibilityChecker{}
}

func (c *compatibilityChecker) Check(parts []model.Part) error {
	var (
		hull   *model.HullProperties
		engine *model.EngineProperties
		shield *model.ShieldProperties
		weapon *model.WeaponProperties
	)

	// Собираем свойства по слотам: у каждой детали non-nil ровно одно свойство.
	for _, part := range parts {
		props := part.Properties()

		switch {
		case props.Hull() != nil:
			hull = props.Hull()
		case props.Engine() != nil:
			engine = props.Engine()
		case props.Shield() != nil:
			shield = props.Shield()
		case props.Weapon() != nil:
			weapon = props.Weapon()
		}
	}

	// Правило 1: корпус должен выдерживать нагрузку двигателя.
	if hull != nil && engine != nil && !hull.CanSupport(engine) {
		return errs.ErrIncompatibleParts
	}

	// Правило 2: плазменный щит несовместим с лазерным оружием.
	if shield != nil && weapon != nil && shield.ConfilctsWith(weapon) {
		return errs.ErrIncompatibleParts
	}

	return nil
}
