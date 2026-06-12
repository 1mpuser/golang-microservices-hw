package converter

import (
	"github.com/google/uuid"

	"github.com/1mpuser/order/internal/model"
	orderv1 "github.com/1mpuser/shared/pkg/openapi/order/v1"
	paymentv1 "github.com/1mpuser/shared/pkg/proto/payment/v1"
)

type PayDto struct {
	TransactionUUID uuid.UUID
}

type CreateModelDto struct {
	OrderUUID  uuid.UUID
	TotalPrice int64
}

func PayModelToDto(transactionUUID uuid.UUID) *PayDto {
	return &PayDto{
		TransactionUUID: transactionUUID,
	}
}

func PaymentMethodToModel(paymentMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return ""
	}
}

// PaymentMethodFromOpenAPI преобразует способ оплаты из OpenAPI в proto-представление.
func PaymentMethodFromOpenAPI(paymentMethod orderv1.PaymentMethod) paymentv1.PaymentMethod {
	switch paymentMethod {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// OrderToDTO преобразует доменную модель заказа в DTO для HTTP-ответа.
func OrderToDTO(order model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  order.OrderUUID,
		HullUUID:   order.HullUUID,
		EngineUUID: order.EngineUUID,
		TotalPrice: order.TotalPrice,
		Status:     orderv1.OrderStatus(order.Status),
		CreatedAt:  order.CreatedAt,
	}

	if order.ShieldUUID != nil {
		dto.ShieldUUID.SetTo(*order.ShieldUUID)
	}

	if order.WeaponUUID != nil {
		dto.WeaponUUID.SetTo(*order.WeaponUUID)
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID.SetTo(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		dto.PaymentMethod.SetTo(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return dto
}
