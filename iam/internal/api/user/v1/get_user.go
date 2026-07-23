package v1

import (
	"context"

	"github.com/1mpuser/iam/internal/api/converter"
	userv1 "github.com/1mpuser/shared/pkg/proto/user/v1"
)

// Реализовать (неделя 6): хендлер GetUser: user_uuid → service.GetUser → GetUserResponse{user} (converter).
func (a *api) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := a.iAmService.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
