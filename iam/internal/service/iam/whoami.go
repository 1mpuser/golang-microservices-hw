package iam

import (
	"context"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/model"
	"github.com/1mpuser/iam/internal/repository/converter"
)

// Реализовать (неделя 6): Whoami (контракт).
//
//	пустой session_uuid → ErrEmptySessionID
//	SessionRepository.Get → нет/истекла → ErrSessionNotFound
//	вернуть session + user (UserRepository.GetByUUID)
func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	if len(sessionUUID) == 0 {
		return model.Session{}, model.User{}, errs.ErrEmptySessionID
	}

	session, err := s.sessionRepo.Get(ctx, sessionUUID)
	if err != nil {
		return model.Session{}, model.User{}, errs.ErrSessionNotFound
	}

	sessionModel := converter.SessionRedisViewToModel(*session)

	user, err := s.GetUser(ctx, session.UserUUID)
	if err != nil {
		return model.Session{}, model.User{}, errs.ErrSessionNotFound
	}

	return sessionModel, user, nil
}
