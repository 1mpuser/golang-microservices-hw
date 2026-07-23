package v1

import (
	"context"

	"github.com/1mpuser/iam/internal/service/input"
	userv1 "github.com/1mpuser/shared/pkg/proto/user/v1"
)

// Реализовать (неделя 6): хендлер Register: RegisterRequest → input.RegisterInput → service.Register → RegisterResponse{user_uuid}.
func (a *api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	info := req.GetInfo()

	userUUID, err := a.iAmService.Register(ctx, input.RegisterInput{
		Login:    info.Info.GetLogin(),
		Password: info.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
