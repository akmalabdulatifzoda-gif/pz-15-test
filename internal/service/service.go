package service

import (
	"fmt"
	"strings"
)

type Service struct{ repo UserRepo }

func New(repo UserRepo) *Service { return &Service{repo: repo} }

// NormalizeEmail приводит email к виду, в котором он хранится: без пробелов по краям,
// в нижнем регистре. Возвращает ErrInvalidEmail, если адрес явно некорректный.
func NormalizeEmail(email string) (string, error) {
	e := strings.ToLower(strings.TrimSpace(email))
	local, domain, ok := strings.Cut(e, "@")
	if !ok || local == "" || domain == "" || strings.Contains(domain, "@") ||
		!strings.Contains(domain, ".") || strings.ContainsAny(e, " \t") {
		return "", fmt.Errorf("%w: %q", ErrInvalidEmail, email)
	}
	return e, nil
}

// FindIDByEmail возвращает ID пользователя по email или 0 и ошибку.
//   - некорректный email -> ErrInvalidEmail, репозиторий не вызывается;
//   - пользователя нет -> ErrNotFound (как вернул репозиторий);
//   - прочие ошибки репозитория оборачиваются с контекстом, но errors.Is по ним работает.
func (s *Service) FindIDByEmail(email string) (int64, error) {
	e, err := NormalizeEmail(email)
	if err != nil {
		return 0, err
	}
	u, err := s.repo.ByEmail(e)
	if err != nil {
		return 0, fmt.Errorf("find user by email: %w", err)
	}
	return u.ID, nil
}
