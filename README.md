# Сервис сокращения URL
![coverage](https://raw.githubusercontent.com/real-splendid/url-shortener-practicum/badges/.badges/18/merge/coverage.svg)

## Запуск автотестов
```
make test-iteration
make test
```

## Запуск линтера
```
make vet
```

## Запуск бенчмарков для профилирования
`go test ./internal/handlers/... -bench='.' -memprofile=./profiles/base.pprof`
