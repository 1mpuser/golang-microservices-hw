package record

import (
	"time"
)

// Реализовать (неделя 6): DB-представление строки таблицы users (скалярные поля + теги db:"...").
// Пара record + converter — как в inventory/internal/repository/record.
//
//	uuid, login, password_hash, created_at, updated_at
type User struct {
	UUID         string
	Login        string
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
