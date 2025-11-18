# Go Numerator Service

Микросервис на Go, предоставляющий REST API для работы с числами. Принимает числа через POST-запрос, сохраняет в
PostgreSQL и возвращает отсортированный список.

## 🚀 Быстрый старт

### Запуск через Docker Compose

```bash
# Клонирование репозитория
git clone https://github.com/UdinSemen/golang-test-task.git
cd golang-test-task

# Запуск сервиса
docker-compose up --build
```

Сервис будет доступен на http://localhost:8082
Локальный запуск

1) Установите зависимости:

     ```bash
     go mod tidy
     ```

2) Создайте .env файл

    ```bash
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=postgres
    ```
3) Запустите приложение 
   ```bash
   docker-compose up --build
   ```

📡 API
Добавить число
```bash
POST /v1/numerator/num
Content-Type: application/json

{
  "num": 42
}
```
Ответ:
```json
{
  "success": true,
  "body": {
    "nums": [1, 2, 42]
  }
}

```

📁 Структура проекта
```
├── cmd/app/main.go          # Точка входа
├── internal/
│   ├── adapters/http/       # HTTP модели
│   ├── domain/usecase/      # Бизнес-логика
│   ├── infrastructure/     # Реализации (репозитории, роутеры)
│   └── config/              # Конфигурация
├── deployments/
│   ├── docker-compose.yaml   # Docker Compose
│   └── service/config.yaml  # Конфиг сервиса
├── migrations/              # Миграции БД
└── README.md
```

🧪 Тестирование
### Unit-тесты
```bash
go test -v ./...

# Запуск с coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
### Интеграционные тесты
```bash
go test -v -tags=integration ./...

# Запуск с coverage
go test -coverprofile=coverage.out -tags=integration ./...
go tool cover -html=coverage.out
``` 
