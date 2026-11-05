package service

import "errors"

// User — пользователь системы.
type User struct {
	ID    int64
	Email string
}

var (
	// ErrNotFound — пользователь не найден (возвращает репозиторий).
	ErrNotFound = errors.New("not found")
	// ErrInvalidEmail — email не прошёл проверку (возвращает сервис, до обращения к репозиторию).
	ErrInvalidEmail = errors.New("invalid email")
)

// UserRepo — зависимость сервиса. В продакшене — PostgreSQL,
// в unit-тестах — стаб или мок, без базы данных.
type UserRepo interface {
	ByEmail(email string) (User, error)
}
