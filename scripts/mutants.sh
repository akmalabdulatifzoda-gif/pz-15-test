#!/usr/bin/env bash
# Простая мутационная проверка: по очереди вносим в КОПИЮ проекта типичную ошибку
# и смотрим, какие тесты её ловят. Мутант, который «выжил», означает пробел в тестах
# (или эквивалентный мутант — изменение, не меняющее поведение).
set -uo pipefail
cd "$(dirname "$0")/.."
TMP=$(mktemp -d); trap 'rm -rf "$TMP"' EXIT

mutant() { # имя файл "было" "стало"
  rm -rf "$TMP/p"; cp -r . "$TMP/p"
  python3 - "$TMP/p/$2" "$3" "$4" <<'PY'
import sys; p,a,b=sys.argv[1:]; s=open(p).read(); assert a in s, a; open(p,"w").write(s.replace(a,b,1))
PY
  echo "# Мутант: $1  ($2: «$3» → «$4»)"
  if (cd "$TMP/p" && go test ./... 2>&1 | grep -E -- "--- FAIL|^(ok|FAIL)\s"); then :; fi
  echo
}

mutant "граница в Clip"               internal/stringsx/stringsx.go 'if len(s) <= max {' 'if len(s) < max {'
mutant "нет проверки деления на ноль" internal/mathx/mathx.go       'if b == 0 {' 'if b == 0 && a != 0 {'
mutant "сервис не нормализует email"  internal/service/service.go   'u, err := s.repo.ByEmail(e)' 'u, err := s.repo.ByEmail(email); _ = e'
mutant "ошибка репозитория без %w"    internal/service/service.go   '"find user by email: %w"' '"find user by email: %v"'
mutant "max<0 не обрабатывается"      internal/stringsx/stringsx.go $'\tif max < 0 {\n\t\tmax = 0\n\t}\n' ''
