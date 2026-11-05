package mathx

import (
	"math"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- «Чистый» testing ---

func TestSum_Table(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{2, 3, 5},
		{10, -5, 5},
		{-4, -6, -10},
		{0, 0, 0},
		{math.MaxInt, 0, math.MaxInt},
		{math.MinInt, 0, math.MinInt},
	}
	for _, c := range cases {
		if got := Sum(c.a, c.b); got != c.want {
			t.Errorf("Sum(%d,%d) = %d; want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestDivide_OkAndError(t *testing.T) {
	got, err := Divide(10, 2)
	if err != nil {
		t.Fatalf("Divide(10,2): неожиданная ошибка %v", err)
	}
	if got != 5 {
		t.Errorf("Divide(10,2) = %d; want 5", got)
	}

	got, err = Divide(10, 0)
	if err == nil {
		t.Fatal("Divide(10,0): ожидалась ошибка, получено nil")
	}
	if got != 0 {
		t.Errorf("Divide(10,0) при ошибке вернул %d; want 0", got)
	}
}

// --- testify + t.Run ---

func TestDivide_Table(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr error
	}{
		{"обычное деление", 6, 3, 2, nil},
		{"дробная часть отбрасывается", 7, 2, 3, nil},
		{"отрицательное делимое — округление к нулю", -7, 2, -3, nil},
		{"отрицательный делитель", 7, -2, -3, nil},
		{"ноль делится на число", 0, 5, 0, nil},
		{"деление на ноль", 6, 0, 0, ErrDivideByZero},
		{"ноль на ноль", 0, 0, 0, ErrDivideByZero},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, got, "при ошибке результат должен быть нулевым")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- Паники ---

func TestMustDivide(t *testing.T) {
	assert.Equal(t, 4, MustDivide(8, 2))
	require.PanicsWithError(t, ErrDivideByZero.Error(), func() { MustDivide(1, 0) })
}

// --- Property-based: свойства должны выполняться для любых входных данных ---

func TestSum_Properties(t *testing.T) {
	commutative := func(a, b int) bool { return Sum(a, b) == Sum(b, a) }
	identity := func(a int) bool { return Sum(a, 0) == a }
	require.NoError(t, quick.Check(commutative, nil), "Sum должна быть коммутативной")
	require.NoError(t, quick.Check(identity, nil), "0 — нейтральный элемент")
}

func TestDivide_Property(t *testing.T) {
	// Для любого b != 0: a == q*b + r и |r| < |b|.
	prop := func(a, b int) bool {
		if b == 0 {
			_, err := Divide(a, b)
			return err != nil
		}
		if a == math.MinInt && b == -1 { // переполнение int — вне контракта функции
			return true
		}
		q, err := Divide(a, b)
		r := a - q*b
		return err == nil && a == q*b+r && abs(r) < abs(b)
	}
	require.NoError(t, quick.Check(prop, &quick.Config{MaxCount: 5000}))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// --- Бенчмарки ---

// sink и переменные-аргументы не дают компилятору встроить вызов и выбросить
// неиспользуемый результат: иначе бенчмарк измеряет пустой цикл (~0.1 ns/op).
var (
	sink   int
	x, y   = 123, 456
	da, db = 1000, 7
)

// BenchmarkSum_Naive — вариант из задания: результат отбрасывается, вызов вырезается оптимизатором.
func BenchmarkSum_Naive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum(123, 456)
	}
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = Sum(x, i)
	}
}

func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink, _ = Divide(da+i, db)
	}
}
