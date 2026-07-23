package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/1mpuser/iam/internal/model"
	commonv1 "github.com/1mpuser/shared/pkg/proto/common/v1"
)

// Реализовать (неделя 6): конвертация model.User/Session → common.v1 (User, UserInfo, Session) для gRPC-ответов.
func SessionToProto(s model.Session) *commonv1.Session {
	return &commonv1.Session{
		Uuid:      s.UUID.String(),
		CreatedAt: timestamppb.New(s.CreatedAt),
		ExpiresAt: timestamppb.New(s.ExpiresAt),
	}
}

func UserToProto(u model.User) *commonv1.User {
	user := &commonv1.User{
		Uuid:      u.UUID.String(),
		Info:      &commonv1.UserInfo{Login: u.Login},
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
	if u.UpdatedAt != nil { // model.UpdatedAt — *time.Time, может быть nil
		user.UpdatedAt = timestamppb.New(*u.UpdatedAt)
	}
	return user
}
