// Package app — тонкая обёртка над internal/app для API-тестов.
// Собирает зависимости PaymentService так же, как internal/app/di.go.
package app

import (
	"google.golang.org/grpc"

	paymentAPIv1 "github.com/1mpuser/payment/internal/api/payment/v1"
	paymentService "github.com/1mpuser/payment/internal/service/payment"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

// Interceptors возвращает gRPC-опции сервера.
func Interceptors() []grpc.ServerOption {
	return nil
}

// RegisterServices собирает зависимости PaymentService и регистрирует их на gRPC-сервере.
func RegisterServices(server *grpc.Server) {
	svc := paymentService.NewService()
	api := paymentAPIv1.NewApi(svc)

	paymentv1.RegisterPaymentServiceServer(server, api)
}
