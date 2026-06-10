package v1

import paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"

type api struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService PaymentService
}

func NewApi(paymentService PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
