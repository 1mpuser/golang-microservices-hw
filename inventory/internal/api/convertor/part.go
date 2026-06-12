package convertor

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/1mpuser/inventory/internal/model"
	inventoryv1 "github.com/1mpuser/shared/pkg/proto/inventory/v1"
)

func PartToDto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(*part.CreatedAt),
	}
}
