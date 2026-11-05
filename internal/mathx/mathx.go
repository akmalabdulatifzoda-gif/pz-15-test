// Package mathx — простые арифметические функции для отработки unit-тестов.
package mathx

import "errors"

// ErrDivideByZero возвращается при делении на ноль.
var ErrDivideByZero = errors.New("divide by zero")

// Sum возвращает сумму a и b.
func Sum(a, b int) int { return a + b }

// Divide возвращает целую часть a/b (с отбрасыванием дробной части к нулю).
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

// MustDivide — вариант Divide для заведомо корректных входных данных:
// при делении на ноль паникует (как regexp.MustCompile при неверном шаблоне).
func MustDivide(a, b int) int {
	q, err := Divide(a, b)
	if err != nil {
		panic(err)
	}
	return q
}
