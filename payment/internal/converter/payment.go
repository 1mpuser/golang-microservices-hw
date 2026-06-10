package converter

import (
	"github.com/1mpuser/payment/internal/model"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

func DtoToModel(dto *paymentv1.PayOrderRequest) model.PayRequest {
	return model.PayRequest{
		OrderUUID:     dto.GetOrderUuid(),
		PaymentMethod: paymentMethodToModel(dto.GetPaymentMethod()),
	}
}

func paymentMethodToModel(method paymentv1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
