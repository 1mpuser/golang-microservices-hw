package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/1mpuser/platform/pkg/auth"
)

type IAMClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

// Реализовать (неделя 6): HTTP auth-middleware OrderService (контракт «HTTP Middleware»).
//
//	читает "Authorization: Bearer <session_uuid>"
//	нет заголовка → 401 "отсутствует заголовок Authorization"
//	не Bearer     → 401 "неверный формат Authorization"
//	AuthService.Whoami(session_uuid); ошибка → 401 "недействительная сессия"
//	успех → auth.WithUserUUID + auth.WithSessionUUID (platform/pkg/auth) → next
//	Защищает все эндпоинты /api/v1/orders*.
func NewAuth(iAmClient IAMClient) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "отсутствует заголовок Authorization", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "неверный формат Authorization", http.StatusUnauthorized)
				return
			}

			sessionUid := strings.TrimPrefix(authHeader, "Bearer ")

			userUUID, err := iAmClient.Whoami(r.Context(), sessionUid)
			if err != nil {
				http.Error(w, "недействительная сессия", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			ctx = auth.WithSessionUUID(ctx, sessionUid)
			ctx = auth.WithUserUUID(ctx, userUUID)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
