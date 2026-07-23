package v1

import (
	"context"

	"github.com/1mpuser/iam/internal/api/converter"
	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
)

// Реализовать (неделя 6): хендлер Whoami: session_uuid → service.Whoami → WhoamiResponse{session, user} (converter).
func (a *api) Whoami(ctx context.Context, whoami *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	session, user, err := a.iAmService.Whoami(ctx, whoami.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &authv1.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
