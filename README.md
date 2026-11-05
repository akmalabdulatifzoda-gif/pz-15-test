# ПЗ №15 — Unit-тестирование функций (testing, testify)

**Студент:** Абдулатифзода Акмал, группа ЭФМО-0126
**Дисциплина:** Технологии индустриального программирования (преподаватель Адышкин С. С.)

![tests](https://github.com/akmalabdulatifzoda-gif/pz-15-test/actions/workflows/test.yml/badge.svg)

## Структура

```
internal/
├── mathx/      Sum, Divide, MustDivide          — табличные тесты, ошибки, паники, property-based
├── stringsx/   Clip, ClipRunes (UTF-8)          — t.Run, граничные случаи, fuzz-тест
└── service/    UserRepo (интерфейс), Service    — стаб, ручной мок (spy), testify/mock
```

## Как запустить

```bash
go test ./...                       # все тесты
go test -v ./internal/...           # подробный вывод
go test -cover ./...                # процент покрытия
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out    # html-отчёт покрытия в браузере
go tool cover -func=coverage.out    # покрытие по функциям
go test -run='^$' -bench=. -benchmem ./internal/...              # бенчмарки
go test -run='^$' -fuzz=FuzzClipRunes -fuzztime=10s ./internal/stringsx   # fuzz
```

То же через Makefile: `make test`, `make cover`, `make cover-html`, `make bench`, `make fuzz`.
Тесты запускаются в GitHub Actions на каждый push (`.github/workflows/test.yml`, с `-race`).

## Результаты

| Пакет | Покрытие |
|---|---|
| internal/mathx    | 100.0% |
| internal/service  | 100.0% |
| internal/stringsx | 100.0% |

HTML-отчёт покрытия: [report/coverage.html](report/coverage.html), отчёт — [report/](report).

Что проверено помимо покрытия:

- **Мутационная проверка** (`results/mutants.txt`): в код по очереди вносились 5 ошибок —
  4 из 5 тесты поймали; выживший мутант (`<=` → `<` в Clip) эквивалентен исходному коду.
- **Property-based** (`testing/quick`): коммутативность Sum, `a == q*b + r, |r| < |b|` для Divide.
- **Fuzz** (`FuzzClipRunes`): ~190 тыс. случайных входов за 10 с, результат всегда валидный UTF-8 и префикс исходной строки.
- **Бенчмарки** (`results/bench.txt`): найдена ловушка «вырезанного» вызова (0.15 ns/op) и исправлена через sink;
  ClipRunes ускорена в 13 раз (1 690 → 122 ns/op) за счёт одного прохода вместо подсчёта всех символов.
