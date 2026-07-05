package model

import (
	"fmt"

	errs "github.com/1mpuser/inventory/internal/errors"
)

type WeaponType string

const (
	WeaponTypeLaser   WeaponType = "laser"
	WeaponTypeMissile WeaponType = "missile"
)

type WeaponProperties struct {
	weaponType WeaponType
}

func (w *WeaponProperties) Weapon() WeaponType { return w.weaponType }

func NewWeaponPropeties(weaponType string) (PartProperties, error) {
	tmpWeaponType := WeaponType(weaponType)

	switch tmpWeaponType {
	case WeaponTypeLaser, WeaponTypeMissile:
		return PartProperties{
			weapon: &WeaponProperties{
				weaponType: tmpWeaponType,
			},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип оружия: %q: %w", weaponType, errs.ErrIncompatibleParts)
	}
}
