package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/service/domain"
)

func hullPart(t *testing.T, strength int) model.Part {
	t.Helper()
	props, err := model.NewHullProperties(strength)
	require.NoError(t, err)
	return model.RestorePart("hull", "Корпус", "", model.PartTypeHull, 0, 1, 0, props, time.Now())
}

func enginePart(t *testing.T, requiredStrength int) model.Part {
	t.Helper()
	props, err := model.NewEngineProperties("B", requiredStrength)
	require.NoError(t, err)
	return model.RestorePart("engine", "Двигатель", "", model.PartTypeEngine, 0, 1, 0, props, time.Now())
}

func shieldPart(t *testing.T, shieldType string) model.Part {
	t.Helper()
	props, err := model.NewShieldPropeties(shieldType)
	require.NoError(t, err)
	return model.RestorePart("shield", "Щит", "", model.PartTypeShield, 0, 1, 0, props, time.Now())
}

func weaponPart(t *testing.T, weaponType string) model.Part {
	t.Helper()
	props, err := model.NewWeaponPropeties(weaponType)
	require.NoError(t, err)
	return model.RestorePart("weapon", "Оружие", "", model.PartTypeWeapon, 0, 1, 0, props, time.Now())
}

func TestCompatibilityChecker_Check(t *testing.T) {
	t.Parallel()

	checker := domain.NewCompatibilityChecker()

	tests := []struct {
		name    string
		parts   []model.Part
		wantErr error
	}{
		{
			name:    "корпус выдерживает двигатель — совместимо",
			parts:   []model.Part{hullPart(t, 100), enginePart(t, 70)},
			wantErr: nil,
		},
		{
			name:    "слабый корпус — несовместимо",
			parts:   []model.Part{hullPart(t, 50), enginePart(t, 70)},
			wantErr: errs.ErrIncompatibleParts,
		},
		{
			name:    "плазменный щит + лазер — несовместимо",
			parts:   []model.Part{hullPart(t, 150), enginePart(t, 70), shieldPart(t, "plasma"), weaponPart(t, "laser")},
			wantErr: errs.ErrIncompatibleParts,
		},
		{
			name:    "энергощит + лазер — совместимо",
			parts:   []model.Part{hullPart(t, 150), enginePart(t, 70), shieldPart(t, "energy"), weaponPart(t, "laser")},
			wantErr: nil,
		},
		{
			name:    "плазменный щит + ракета — совместимо",
			parts:   []model.Part{hullPart(t, 150), enginePart(t, 70), shieldPart(t, "plasma"), weaponPart(t, "missile")},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := checker.Check(tc.parts)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
