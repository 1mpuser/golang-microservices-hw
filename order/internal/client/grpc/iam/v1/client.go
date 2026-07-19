package v1

// Реализовать (неделя 6): gRPC-клиент IAMService для OrderService (обёртка над authv1.AuthServiceClient).
//   Whoami(ctx, sessionUUID) → user_uuid/login; используется HTTP middleware.
//   По аналогии с order/internal/client/grpc/inventory/v1/client.go.
