package iam

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/model"
	"github.com/1mpuser/iam/internal/repository/converter"
)

func (s *service) GetUser(ctx context.Context, uid string) (model.User, error) {
	_, err := uuid.Parse(uid)
	if err != nil {
		return model.User{}, errs.ErrInvalidUUID
	}

	record, err := s.userRepo.GetByUuid(ctx, uid)
	if err != nil {
		return model.User{}, errs.ErrUserNotFound
	}

	return converter.RecordToModel(*record), nil
}
