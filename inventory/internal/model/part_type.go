package model

import (
	"fmt"

	errs "github.com/1mpuser/inventory/internal/errors"
)

type PartType string

const (
	PartTypeUnspecified PartType = ""
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

func NewPartType(s string) (PartType, error) {
	partType := PartType(s)

	switch partType {
	case PartTypeHull, PartTypeEngine, PartTypeShield, PartTypeWeapon:
		return partType, nil
	default:
		return "", fmt.Errorf("неизвестный формат детали: %q %w", s, errs.ErrInvalidProperties)
	}
}
