package stringsx

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestClip(t *testing.T) {
	tests := []struct {
		name string
		s    string
		max  int
		want string
	}{
		{"пустая строка", "", 5, ""},
		{"пустая строка и max=0", "", 0, ""},
		{"max=0", "hello", 0, ""},
		{"max<0 считается нулём", "hello", -3, ""},
		{"max == len(s)", "hello", 5, "hello"},
		{"max > len(s)", "hello", 10, "hello"},
		{"обычная обрезка", "hello", 3, "hel"},
		{"max=1", "hello", 1, "h"},
		{"пробелы сохраняются", "  ab  ", 3, "  a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clip(tt.s, tt.max)
			assert.Equal(t, tt.want, got)
			assert.LessOrEqual(t, len(got), max(tt.max, 0), "длина результата не превышает max")
		})
	}
}

// Документирует известное ограничение Clip: обрезка по байтам ломает UTF-8.
func TestClip_BreaksUTF8(t *testing.T) {
	got := Clip("Привет", 3) // «П» и половина «р» (по 2 байта на символ)
	assert.Len(t, got, 3)
	assert.False(t, utf8.ValidString(got), "Clip режет по байтам — строка становится невалидной")
}

func TestClipRunes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		max  int
		want string
	}{
		{"кириллица", "Привет", 3, "При"},
		{"смешанный текст", "Go-разработка", 4, "Go-р"},
		{"эмодзи — один символ", "😀😀😀", 2, "😀😀"},
		{"max больше длины", "Мир", 10, "Мир"},
		{"max равен длине в рунах", "Мир", 3, "Мир"},
		{"max=0", "Мир", 0, ""},
		{"max<0", "Мир", -1, ""},
		{"пустая строка", "", 3, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClipRunes(tt.s, tt.max)
			assert.Equal(t, tt.want, got)
			assert.True(t, utf8.ValidString(got), "результат должен оставаться валидным UTF-8")
		})
	}
}

// Fuzz-тест: случайные строки и длины, проверяются инварианты.
// Запуск: go test -fuzz=FuzzClipRunes -fuzztime=10s ./internal/stringsx
func FuzzClipRunes(f *testing.F) {
	for _, seed := range []string{"", "hello", "Привет", "😀a😀", strings.Repeat("я", 50)} {
		f.Add(seed, 3)
	}
	f.Fuzz(func(t *testing.T, s string, n int) {
		if !utf8.ValidString(s) {
			t.Skip("невалидный UTF-8 на входе вне контракта")
		}
		got := ClipRunes(s, n)
		if !utf8.ValidString(got) {
			t.Fatalf("ClipRunes(%q, %d) = %q — невалидный UTF-8", s, n, got)
		}
		if !strings.HasPrefix(s, got) {
			t.Fatalf("результат %q не является префиксом %q", got, s)
		}
		if limit := max(n, 0); utf8.RuneCountInString(got) > limit {
			t.Fatalf("в результате %d рун, а max=%d", utf8.RuneCountInString(got), n)
		}
	})
}

var sink string

func BenchmarkClip(b *testing.B) {
	s := strings.Repeat("hello ", 100)
	for i := 0; i < b.N; i++ {
		sink = Clip(s, 50+i%2)
	}
}

// Длинная строка (700 символов), обрезка до 50 — частый случай для превью текста.
func BenchmarkClipRunes(b *testing.B) {
	s := strings.Repeat("привет ", 100)
	for i := 0; i < b.N; i++ {
		sink = ClipRunes(s, 50)
	}
}

// Короткая строка, обрезка не нужна.
func BenchmarkClipRunes_Short(b *testing.B) {
	s := "Привет, мир"
	for i := 0; i < b.N; i++ {
		sink = ClipRunes(s, 50)
	}
}
