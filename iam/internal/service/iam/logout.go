package iam

import (
	"context"

	errs "github.com/1mpuser/iam/internal/errors"
)

// Реализовать (неделя 6): Logout (контракт).
//
//	пустой session_uuid → ErrEmptySessionID
//	SessionRepository.Delete — идемпотентно (нет сессии → OK)
func (s service) Logout(ctx context.Context, sessionUid string) error {
	if len(sessionUid) == 0 {
		return errs.ErrEmptySessionID
	}

	return s.sessionRepo.Delete(ctx, sessionUid)
}
