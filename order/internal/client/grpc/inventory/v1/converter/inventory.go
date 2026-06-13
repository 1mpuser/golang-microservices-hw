package converter

import (
	"github.com/1mpuser/order/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func PartTypeFromProto(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return ""
	}
}

func DTOToModel(dto *inventoryv1.ListPartsResponse) []model.Part {
	parts := make([]model.Part, 0, len(dto.Parts))

	for _, part := range dto.Parts {
		parts = append(parts, model.Part{
			UUID:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			PartType:      PartTypeFromProto(part.PartType),
			StockQuantity: part.StockQuantity,
			CreatedAt:     new(part.CreatedAt.AsTime()),
		})
	}

	return parts
}
