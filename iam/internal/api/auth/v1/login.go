package v1

import (
	"context"

	"github.com/1mpuser/iam/internal/service/input"
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
)

// Реализовать (неделя 6): хендлер Login: authv1.LoginRequest → input.LoginInput → service.Login → LoginResponse{session_uuid}.
func (a *api) Login(ctx context.Context, login *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	sessionUUID, err := a.iAmService.Login(ctx, input.LoginInput{
		Login:    login.GetLogin(),
		Password: login.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
