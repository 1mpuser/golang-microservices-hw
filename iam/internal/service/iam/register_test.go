package iam_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/iam/internal/errors"
	iamService "github.com/1mpuser/iam/internal/service/iam"
	"github.com/1mpuser/iam/internal/service/iam/mocks"
	"github.com/1mpuser/iam/internal/service/input"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		in        input.RegisterInput
		setupMock func(userRepo *mocks.UserRepository)
		wantErr   error
	}{
		{
			name: "успешная регистрация",
			in:   input.RegisterInput{Login: "user", Password: "password123"},
			setupMock: func(userRepo *mocks.UserRepository) {
				userRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "пустой логин",
			in:        input.RegisterInput{Login: "", Password: "password123"},
			setupMock: func(_ *mocks.UserRepository) {},
			wantErr:   errs.ErrInvalidLogin,
		},
		{
			name:      "слишком короткий пароль",
			in:        input.RegisterInput{Login: "user", Password: "short"},
			setupMock: func(_ *mocks.UserRepository) {},
			wantErr:   errs.ErrWeakPassword,
		},
		{
			name: "логин уже занят",
			in:   input.RegisterInput{Login: "user", Password: "password123"},
			setupMock: func(userRepo *mocks.UserRepository) {
				userRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errs.ErrUserAlreadyExists)
			},
			wantErr: errs.ErrUserAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userRepo := mocks.NewUserRepository(t)
			sessionRepo := mocks.NewSessionRepository(t) // в Register не участвует
			tc.setupMock(userRepo)

			svc := iamService.NewService(userRepo, sessionRepo)
			gotUUID, err := svc.Register(ctx, tc.in)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, gotUUID)
				return
			}

			require.NoError(t, err)
			_, parseErr := uuid.Parse(gotUUID)
			assert.NoError(t, parseErr, "Register должен вернуть валидный user_uuid")
		})
	}
}
