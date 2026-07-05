package app

import (
	"context"

	paymentAPIv1 "github.com/1mpuser/payment/internal/api/payment/v1"
	paymentService "github.com/1mpuser/payment/internal/service/payment"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentSvc     paymentAPIv1.PaymentService
	paymentHandler paymentv1.PaymentServiceServer
}

// PaymentService возвращает сервис бизнес-логики платежей.
func (d *diContainer) PaymentService(_ context.Context) paymentAPIv1.PaymentService {
	if d.paymentSvc == nil {
		d.paymentSvc = paymentService.NewService()
	}
	return d.paymentSvc
}

// PaymentAPI возвращает gRPC-обработчик PaymentService.
func (d *diContainer) PaymentAPI(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentHandler == nil {
		d.paymentHandler = paymentAPIv1.NewApi(d.PaymentService(ctx))
	}
	return d.paymentHandler
}
