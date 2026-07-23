package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/1mpuser/platform/pkg/auth"
	kafkamw "github.com/1mpuser/platform/pkg/middleware/kafka"
)

// Реализовать (неделя 6): gRPC server (incoming) interceptor InventoryService (контракт «gRPC Interceptor»).
//
//	metadata.FromIncomingContext → нет metadata → Unauthenticated "отсутствует metadata"
//	ключ "session-uuid" отсутствует → Unauthenticated "отсутствует session-uuid"
//	AuthService.Whoami(session_uuid); ошибка → Unauthenticated "недействительная сессия"
//	успех → auth.WithUserUUID → handler. Защищает все 6 методов Inventory.
type IAMClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

func NewAuth(iAmClient IAMClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "отсутствует metadata")
		}

		values := md.Get(kafkamw.SessionHeaderKey)

		if len(values) == 0 || values[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "отсутствует session-uuid")
		}

		userUUID, err := iAmClient.Whoami(ctx, values[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "недействительная сессия")
		}

		ctx = auth.WithUserUUID(ctx, userUUID)

		return handler(ctx, req)

	}
}
