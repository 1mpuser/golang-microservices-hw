package errs

import "errors"

var (
	ErrPartNotFound      = errors.New("деталь не найдена")
	ErrInvalidUUID       = errors.New("неверный формат UUID")
	ErrOutOfStock        = errors.New("деталь отсутствует на складе")
	ErrNothingToRelease  = errors.New("нечего освобождать")
	ErrNothingToCommit   = errors.New("нечего выдавать")
	ErrIncompatibleParts = errors.New("детали несовместимы")
	ErrInvalidProperties = errors.New("некорректные свойства детали")
)
