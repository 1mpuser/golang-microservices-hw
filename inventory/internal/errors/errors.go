package errs

import "errors"

var (
	ErrPartNotFound  = errors.New("деталь не найдена")
	ErrInvalidFormat = errors.New("неверный формат uuid")
)
