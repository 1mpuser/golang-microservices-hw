package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/1mpuser/iam/internal/errors"
)

// ErrorUnaryInterceptor — единый gRPC-перехватчик, который превращает доменные
// ошибки сервиса (errs.*) в gRPC-статусы по таблице контрактов. Хендлеры
// возвращают «чистые» доменные ошибки, а маппинг на codes.* живёт здесь одним
// местом (в отличие от inventory, где каждый хендлер мапит сам).
//
// Подключается в app/di.go:
//
//	grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ErrorUnaryInterceptor))
func ErrorUnaryInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	// Ошибка уже пришла с gRPC-статусом (например, из вложенного вызова) —
	// не переоборачиваем, отдаём как есть.
	if _, ok := status.FromError(err); ok {
		return resp, err
	}

	switch {
	case errors.Is(err, errs.ErrInvalidLogin),
		errors.Is(err, errs.ErrWeakPassword),
		errors.Is(err, errs.ErrEmptyCredential),
		errors.Is(err, errs.ErrEmptySessionID),
		errors.Is(err, errs.ErrInvalidUUID):
		return resp, status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrSessionNotFound):
		return resp, status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errs.ErrUserAlreadyExists):
		return resp, status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errs.ErrUserNotFound):
		return resp, status.Error(codes.NotFound, err.Error())

	default:
		// Неизвестную ошибку наружу детально не раскрываем.
		return resp, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
