package interceptor

import (
	"context"

	"github.com/1mpuser/platform/pkg/auth"
	kafkamw "github.com/1mpuser/platform/pkg/middleware/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Реализовать (неделя 6): gRPC client (outgoing) interceptor OrderService.
//
//	auth.SessionUUIDFromContext(ctx) → metadata.AppendToOutgoingContext(ctx, "session-uuid", uuid)
//	подключается к gRPC-клиенту InventoryService в app/di.go.
//	Без него Inventory не получит токен сессии и упадёт на Unauthenticated.
func SessionForwarder(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if sessionUUID, ok := auth.SessionUUIDFromContext(ctx); ok {
		ctx = metadata.AppendToOutgoingContext(ctx, kafkamw.SessionHeaderKey, sessionUUID)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
