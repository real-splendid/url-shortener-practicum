# Сервис сокращения URL

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
