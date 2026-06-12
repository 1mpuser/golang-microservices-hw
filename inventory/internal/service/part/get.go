package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, id string) (model.Part, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Part{}, errs.ErrInvalidFormat
	}

	return s.partRepository.Get(ctx, uid)
}
