package iam

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/repository/record"
	"github.com/1mpuser/iam/internal/service/input"
)

// Реализовать (неделя 6): Register (контракт).
//
//	валидация: пустой логин → ErrInvalidLogin; пароль < 8 → ErrWeakPassword
//	bcrypt.GenerateFromPassword(DefaultCost) → UserRepository.Create → user_uuid
func (s *service) Register(ctx context.Context, in input.RegisterInput) (string, error) {
	login := in.Login

	if len(login) == 0 {
		return "", errs.ErrInvalidLogin
	}

	password := in.Password

	if len(password) < 8 {
		return "", errs.ErrWeakPassword
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userUid := uuid.New()

	recordUser := record.User{
		UUID:         userUid.String(),
		Login:        login,
		PasswordHash: string(hashPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, recordUser)
	if err != nil {
		return "", err
	}

	return userUid.String(), nil
}
