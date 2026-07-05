package model

import (
	"fmt"

	errs "github.com/1mpuser/inventory/internal/errors"
)

type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

type EngineProperties struct {
	class            EngineClass
	requiredStrength int
}

func (e *EngineProperties) Class() EngineClass    { return e.class }
func (e *EngineProperties) RequiredStrength() int { return e.requiredStrength }

func NewEngineProperties(class string, requiredStrength int) (PartProperties, error) {
	switch EngineClass(class) {
	case EngineClassA, EngineClassB, EngineClassC:
		return PartProperties{engine: &EngineProperties{
			class:            EngineClass(class),
			requiredStrength: requiredStrength,
		}}, nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный класс двигателя: %q: %w", class, errs.ErrIncompatibleParts)
	}
}
