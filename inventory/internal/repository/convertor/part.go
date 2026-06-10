package convertor

import (
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/record"
)

func PartToModel(p record.Part) model.Part {
	return model.Part{
		UUID:          p.UUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      p.PartType,
		StockQuantity: p.StockQuantity,
		CreatedAt:     &p.CreatedAt,
	}
}
