package errs

import "errors"

var (
	ErrOrderNotFound          = errors.New("заказ не найден")
	ErrOrderAlreadyPaid       = errors.New("заказ уже оплачен")
	ErrOrderCancelled         = errors.New("заказ отменён")
	ErrPartNotFound           = errors.New("деталь не найдена")
	ErrOutOfStock             = errors.New("деталь отсутствует на складе")
	ErrInventoryUnavailable   = errors.New("сервис склада недоступен")
	ErrOrderAssembled         = errors.New("заказ уже собран")
	ErrInventoryPartsNotFound = errors.New("на складе не найдены детали")
	ErrIncompatibleParts      = errors.New("детали несовместимы")
	ErrInvalidUUID            = errors.New("неверный формат UUID")
	ErrInvalidPaymentMethod   = errors.New("неверный метод оплаты")
	ErrUnauthorized           = errors.New("пользователь не авторизован")
)
