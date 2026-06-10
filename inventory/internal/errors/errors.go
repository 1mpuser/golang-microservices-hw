package errs

import "errors"

var (
	ErrPartNotFound  = errors.New("Деталь не найдена")
	ErrInvalidFormat = errors.New("Неверный формат uuid")
)
