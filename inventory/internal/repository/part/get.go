package part

import (
	"context"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/convertor"
	"github.com/google/uuid"
)

func (r *repository) Get(_ context.Context, uuid uuid.UUID) (model.Part, error) {

	r.mu.RLock()

	defer r.mu.RUnlock()

	part, ok := r.data[uuid]

	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return convertor.PartToModel(part), nil
}
