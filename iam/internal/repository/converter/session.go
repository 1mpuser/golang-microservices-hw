package converter

import (
	"time"

	"github.com/1mpuser/iam/internal/model"
	redisview "github.com/1mpuser/iam/internal/repository/redis_view"
	"github.com/google/uuid"
)

func SessionRedisViewToModel(r redisview.SessionRedisView) model.Session {
	uidParsed := uuid.MustParse(r.UUID)
	userUidParse := uuid.MustParse(r.UserUUID)
	timeCreatedAt, _ := time.Parse(time.RFC3339, r.CreatedAt)
	timeExpiresAt, _ := time.Parse(time.RFC3339, r.ExpiresAt)

	return model.Session{
		UUID:      uidParsed,
		UserUUID:  userUidParse,
		Login:     r.Login,
		CreatedAt: timeCreatedAt,
		ExpiresAt: timeExpiresAt,
	}
}

func SessionModelToRedisView(m model.Session) redisview.SessionRedisView {
	return redisview.SessionRedisView{
		UUID:      m.UUID.String(),
		UserUUID:  m.UserUUID.String(),
		Login:     m.Login,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		ExpiresAt: m.ExpiresAt.Format(time.RFC3339),
	}
}
