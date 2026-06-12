package converter

import (
	"github.com/1mpuser/order/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func DTOToModel(dto *inventoryv1.ListPartsResponse) []model.Part {
	parts := make([]model.Part, 0, len(dto.Parts))

	for _, part := range dto.Parts {
		parts = append(parts, model.Part{
			UUID:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			PartType:      part.PartType,
			StockQuantity: part.StockQuantity,
			CreatedAt:     new(part.CreatedAt.AsTime()),
		})
	}

	return parts
}
