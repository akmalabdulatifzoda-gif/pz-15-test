// Package stringsx — функции для работы со строками.
package stringsx

// Clip возвращает s, обрезанную до max байт.
// Отрицательный max считается нулём. UTF-8 не учитывается (как в задании):
// для кириллицы обрезка может разрезать символ — для этого есть ClipRunes.
func Clip(s string, max int) string {
	if max < 0 {
		max = 0
	}
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// ClipRunes — UTF-8-безопасная версия: обрезает по символам (рунам), а не по байтам.
// Один проход по строке и остановка на max-м символе: стоимость O(max), а не O(len(s)).
func ClipRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max { // в строке не больше max байт — значит, и символов не больше
		return s
	}
	n := 0
	for i := range s { // i — байтовое смещение начала очередной руны
		if n == max {
			return s[:i]
		}
		n++
	}
	return s
}
