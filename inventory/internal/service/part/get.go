package part

import (
	"context"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, id string) (model.Part, error) {

	uid, err := uuid.Parse(id)

	if err != nil {
		return model.Part{}, errs.ErrInvalidFormat
	}

	return s.partRepository.Get(ctx, uid)

}
