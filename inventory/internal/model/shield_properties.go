package model

import (
	"fmt"

	errs "github.com/1mpuser/inventory/internal/errors"
)

type ShieldType string

const (
	ShieldTypeEnergy ShieldType = "energy"
	ShieldTypePlasma ShieldType = "plasma"
)

type ShieldProperties struct {
	shieldType ShieldType
}

func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

func NewShieldPropeties(shieldType string) (PartProperties, error) {
	switch ShieldType(shieldType) {
	case ShieldTypeEnergy, ShieldTypePlasma:
		return PartProperties{
			shield: &ShieldProperties{
				shieldType: ShieldType(shieldType),
			},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип щита: %q: %w", shieldType, errs.ErrIncompatibleParts)
	}
}

func (s *ShieldProperties) ConfilctsWith(w *WeaponProperties) bool {
	return s.ShieldType() == ShieldTypePlasma && w.Weapon() == WeaponTypeLaser
}
