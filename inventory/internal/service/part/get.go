package part

import (
	"context"
	"errors"

	"github.com/google/uuid"

	errs "github.com/1mpuser/inventory/internal/errors"
	"github.com/1mpuser/inventory/internal/model"
	"github.com/1mpuser/inventory/internal/repository/convertor"
)

func (s *service) Get(ctx context.Context, id string) (model.Part, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.Part{}, errs.ErrInvalidFormat
	}

	part, err := s.partRepository.Get(ctx, uid)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return model.Part{}, errs.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return convertor.PartToModel(part), nil
}
