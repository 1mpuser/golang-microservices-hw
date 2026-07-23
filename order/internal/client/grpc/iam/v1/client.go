package v1

import (
	"context"

	"github.com/google/uuid"

	authv1 "github.com/1mpuser/shared/pkg/proto/auth/v1"
)

// client — обёртка над gRPC-клиентом IAMService (AuthService).
type client struct {
	grpcClient authv1.AuthServiceClient
}

func New(grpcClient authv1.AuthServiceClient) *client {
	return &client{
		grpcClient,
	}
}

// Whoami проверяет сессию через IAM и возвращает UUID пользователя.
// Ошибку IAM (Unauthenticated/InvalidArgument) пробрасывает как есть —
// вызывающий (HTTP middleware) решает, как отдать её наружу (401).
func (c *client) Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error) {
	resp, err := c.grpcClient.Whoami(ctx, &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(resp.GetUser().GetUuid())
}
