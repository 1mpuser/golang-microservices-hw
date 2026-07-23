package iam_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/repository/record"
	iamService "github.com/1mpuser/iam/internal/service/iam"
	"github.com/1mpuser/iam/internal/service/iam/mocks"
	"github.com/1mpuser/iam/internal/service/input"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const (
		login    = "user"
		password = "password123"
	)

	// MinCost — тесты быстрее в ~100 раз (подсказка из ДЗ).
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)

	user := &record.User{
		UUID:         uuid.New().String(),
		Login:        login,
		PasswordHash: string(hash),
	}

	errRedis := errors.New("redis недоступен")

	tests := []struct {
		name      string
		in        input.LoginInput
		setupMock func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository)
		wantErr   error
	}{
		{
			name: "успешный вход",
			in:   input.LoginInput{Login: login, Password: password},
			setupMock: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository) {
				userRepo.EXPECT().GetByLogin(mock.Anything, login).Return(user, nil)
				sessionRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "пустые креды",
			in:        input.LoginInput{Login: "", Password: ""},
			setupMock: func(_ *mocks.UserRepository, _ *mocks.SessionRepository) {},
			wantErr:   errs.ErrEmptyCredential,
		},
		{
			name: "логин не найден → ErrInvalidCredentials",
			in:   input.LoginInput{Login: login, Password: password},
			setupMock: func(userRepo *mocks.UserRepository, _ *mocks.SessionRepository) {
				userRepo.EXPECT().GetByLogin(mock.Anything, login).Return(nil, errs.ErrUserNotFound)
			},
			wantErr: errs.ErrInvalidCredentials,
		},
		{
			name: "неверный пароль → ErrInvalidCredentials",
			in:   input.LoginInput{Login: login, Password: "wrongpassword"},
			setupMock: func(userRepo *mocks.UserRepository, _ *mocks.SessionRepository) {
				userRepo.EXPECT().GetByLogin(mock.Anything, login).Return(user, nil)
			},
			wantErr: errs.ErrInvalidCredentials,
		},
		{
			name: "ошибка Redis при создании сессии",
			in:   input.LoginInput{Login: login, Password: password},
			setupMock: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository) {
				userRepo.EXPECT().GetByLogin(mock.Anything, login).Return(user, nil)
				sessionRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(errRedis)
			},
			wantErr: errRedis,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userRepo := mocks.NewUserRepository(t)
			sessionRepo := mocks.NewSessionRepository(t)
			tc.setupMock(userRepo, sessionRepo)

			svc := iamService.NewService(userRepo, sessionRepo)
			gotSession, gotErr := svc.Login(ctx, tc.in)

			if tc.wantErr != nil {
				require.Error(t, gotErr)
				assert.ErrorIs(t, gotErr, tc.wantErr)
				assert.Empty(t, gotSession)
				return
			}

			require.NoError(t, gotErr)
			_, parseErr := uuid.Parse(gotSession)
			assert.NoError(t, parseErr, "Login должен вернуть валидный session_uuid")
		})
	}
}
