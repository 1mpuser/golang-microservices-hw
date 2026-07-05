package convertor

import (
	"encoding/json"
	"fmt"

	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/record"
)

func PartRecordToModel(rec record.PartRecord) (model.Part, error) {
	var propsRec record.PartPropertiesRecord

	if err := json.Unmarshal(rec.Properties, &propsRec); err != nil {
		return model.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
	}

	props, err := partPropertiesFromRecord(propsRec)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать свойства: %w", err)
	}

	partType, err := model.NewPartType(rec.PartType)
	if err != nil {
		return model.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
	}

	return model.RestorePart(
		rec.UUID,
		rec.Name,
		rec.Description,
		partType,
		rec.Price,
		rec.StockQuantity,
		rec.Reserved,
		props,
		rec.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (model.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return model.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return model.NewEngineProperties(rec.Engine.Class, rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return model.NewShieldPropeties(rec.Shield.ShieldType)
	case rec.Weapon != nil:
		return model.NewWeaponPropeties(rec.Weapon.WeaponType)
	default:
		return model.PartProperties{}, nil
	}
}
