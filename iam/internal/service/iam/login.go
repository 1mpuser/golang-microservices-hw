package iam

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/1mpuser/iam/internal/errors"
	redisview "github.com/1mpuser/iam/internal/repository/redis_view"
	"github.com/1mpuser/iam/internal/service/input"
)

// Реализовать (неделя 6): Login (контракт).
//
//	пустые креды → ErrEmptyCredential
//	UserRepository.GetByLogin + bcrypt.CompareHashAndPassword
//	любое несовпадение → ErrInvalidCredentials (не раскрываем, существует ли логин)
//	SessionRepository.Create (TTL 24h) → session_uuid
func (s *service) Login(ctx context.Context, in input.LoginInput) (string, error) {
	login := in.Login
	password := in.Password

	if len(login) == 0 || len(password) == 0 {
		return "", errs.ErrEmptyCredential
	}

	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return "", errs.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errs.ErrInvalidCredentials
	}

	sessionUid := uuid.New()

	redisView := redisview.SessionRedisView{
		UUID:      sessionUid.String(),
		UserUUID:  user.UUID,
		Login:     user.Login,
		CreatedAt: time.Now().Format(time.RFC3339),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Format(time.RFC3339),
	}

	if err := s.sessionRepo.Create(ctx, redisView); err != nil {
		return "", err
	}

	return sessionUid.String(), nil
}
