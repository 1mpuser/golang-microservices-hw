package model

import (
	"fmt"

	errs "github.com/1mpuser/inventory/internal/errors"
)

type HullProperties struct {
	strength int
}

func (h *HullProperties) Strength() int { return h.strength }

func NewHullProperties(strength int) (PartProperties, error) {
	if strength < 30 || strength > 200 {
		return PartProperties{}, fmt.Errorf("прочность корпуса должна быть от 30 до 200, получено %d: %w", strength, errs.ErrInvalidProperties)
	}

	return PartProperties{
		hull: &HullProperties{strength},
	}, nil
}

func (h *HullProperties) CanSupport(e *EngineProperties) bool {
	return h.Strength() >= e.RequiredStrength()
}
