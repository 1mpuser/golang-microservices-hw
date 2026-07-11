package orderpaid

import (
	"context"

	"github.com/1mpuser/assembly/internal/model"
)

// AssemblyService — бизнес-логика сборки, вызывается из обработчика события OrderPaid.
type AssemblyService interface {
	Assemble(ctx context.Context, event model.OrderPaid) error
}
