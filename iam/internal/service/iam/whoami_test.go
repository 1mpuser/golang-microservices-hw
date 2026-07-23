package iam_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/1mpuser/iam/internal/errors"
	"github.com/1mpuser/iam/internal/repository/record"
	redisview "github.com/1mpuser/iam/internal/repository/redis_view"
	iamService "github.com/1mpuser/iam/internal/service/iam"
	"github.com/1mpuser/iam/internal/service/iam/mocks"
)

func TestWhoami(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	sessionUUID := uuid.New().String()
	userUUID := uuid.New().String()

	view := &redisview.SessionRedisView{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     "user",
		CreatedAt: time.Now().Format(time.RFC3339),
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	user := &record.User{
		UUID:  userUUID,
		Login: "user",
	}

	tests := []struct {
		name      string
		sessionID string
		setupMock func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository)
		wantErr   error
	}{
		{
			name:      "успех",
			sessionID: sessionUUID,
			setupMock: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository) {
				sessionRepo.EXPECT().Get(mock.Anything, sessionUUID).Return(view, nil)
				userRepo.EXPECT().GetByUuid(mock.Anything, userUUID).Return(user, nil)
			},
			wantErr: nil,
		},
		{
			name:      "пустой session_uuid",
			sessionID: "",
			setupMock: func(_ *mocks.UserRepository, _ *mocks.SessionRepository) {},
			wantErr:   errs.ErrEmptySessionID,
		},
		{
			name:      "сессия не найдена/истекла",
			sessionID: sessionUUID,
			setupMock: func(_ *mocks.UserRepository, sessionRepo *mocks.SessionRepository) {
				sessionRepo.EXPECT().Get(mock.Anything, sessionUUID).Return(nil, errs.ErrSessionNotFound)
			},
			wantErr: errs.ErrSessionNotFound,
		},
		{
			name:      "пользователь удалён → ErrSessionNotFound",
			sessionID: sessionUUID,
			setupMock: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository) {
				sessionRepo.EXPECT().Get(mock.Anything, sessionUUID).Return(view, nil)
				userRepo.EXPECT().GetByUuid(mock.Anything, userUUID).Return(nil, errs.ErrUserNotFound)
			},
			wantErr: errs.ErrSessionNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userRepo := mocks.NewUserRepository(t)
			sessionRepo := mocks.NewSessionRepository(t)
			tc.setupMock(userRepo, sessionRepo)

			svc := iamService.NewService(userRepo, sessionRepo)
			session, gotUser, err := svc.Whoami(ctx, tc.sessionID)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, sessionUUID, session.UUID.String())
			assert.Equal(t, userUUID, gotUser.UUID.String())
			assert.Equal(t, "user", gotUser.Login)
		})
	}
}
