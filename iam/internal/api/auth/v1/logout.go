package v1

import (
	"context"

	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
)

// Реализовать (неделя 6): хендлер Logout: session_uuid → service.Logout → LogoutResponse{}.
func (a *api) Logout(ctx context.Context, logout *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	err := a.iAmService.Logout(ctx, logout.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &authv1.LogoutResponse{}, nil
}
