PKGS := ./internal/...

.PHONY: test test-v cover cover-html cover-func bench fuzz lint all

## test: все тесты
test:
	go test ./...

## test-v: подробный вывод
test-v:
	go test -v $(PKGS)

## cover: процент покрытия по пакетам
cover:
	go test -cover ./...

## cover-html: html-отчёт покрытия -> coverage.html (и открыть в браузере)
cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Откройте coverage.html или выполните: go tool cover -html=coverage.out"

## cover-func: покрытие по функциям
cover-func:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## bench: бенчмарки
bench:
	go test -run='^$$' -bench=. -benchmem $(PKGS)

## fuzz: fuzz-тест ClipRunes, 10 секунд
fuzz:
	go test -run='^$$' -fuzz=FuzzClipRunes -fuzztime=10s ./internal/stringsx

lint:
	go vet ./...
	gofmt -l .

all: lint test cover
