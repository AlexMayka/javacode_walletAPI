# **Javacode_walletAPI**

REST-сервис на Go для управления кошельками: пополнение, снятие и получение баланса. Проект выполнен в рамках тестового задания и реализован с учётом высокой нагрузки (1000 RPS) и конкурентной работы с одним кошельком.

___

## **📌 Функциональность**
1. [x] POST /api/v1/wallet — пополнение или снятие средств с кошелька
2. [x] GET /api/v1/wallets/{wallet_uuid} — получение баланса кошелька
3. [x] Конкурентная безопасность через транзакции и FOR UPDATE
4. [x] Документация Swagger
5. [x] Логирование через logrus
6. [x] Покрытие тестами (repository, service, controller)

___

## **🛠 Стек технологий**
| Компонент     | Технология              |
|---------------|-------------------------|
| Язык          | Go 1.24                 |
| БД            | PostgreSQL              |
| web-фреймворк | Gin                     |
| Docker        | Docker + Docker Compose |
| Миграции      | Goose (миграции)        |
| Логирование   | Logrus логирование      | 
| Документация  | Swagger                 |
| Тестирование  | SQLMock + Testify       |
| Нагрузка      | WRK (Lua-скрипты)       | 
________

## **📌 API Эндпоинты**
### 💸 Wallet API 
| Метод | URL         | Описание                                                               |
|-------|------------|------------------------------------------------------------------------|
| `GET` | `/api/v1/wallets/{wallet_uuid}` | Получить текущий баланс по UUID кошелька                               |
| `POST` | `/api/v1/wallet` | Выполнить операцию пополнения или снятия средств с указанного кошелька |



## 🚀 Быстрый старт
___
1. Клонируйте репозиторий
```bash
git clone https://github.com/AlexMayka/javacode_walletAPI.git
cd wallet-api
```

2. Запуск через Docker Compose
```bash
docker compose --env-file config.env up -d
```

Это поднимет:
- PostgreSQL 
- Выполнит миграции с помощью Goose
- Запустит Wallet API

4. Swagger-документация
Открыть в браузере: http://localhost:8080/swagger/index.html

___

## 🧪 Тестирование
```bash
go test ./...
```

Нагрузочное тестирование (необходимо установить wrk):
```bash
wrk -t4 -c10 -d30s -s ./load_tests/post.lua http://localhost:8080
```

Результаты тестов находятся в папке load_tests/:
- get_results.txt
- post_results.txt
- mixed_results.txt

___

## ⚙️ Переменные окружения

Файл .env (config.env):
```dotenv
DB_HOST=db
DB_USER=admin
DB_PASSWORD=password
DB_NAME=wallet
DB_PORT=5432
DRIVER=postgres

SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=release
```
___

## 🧩 Архитектура
```bash
* cmd/ — точка входа
* config/ — загрузка конфигурации
* internal/
  * controllers/ — HTTP-обработчики
  * service/ — бизнес-логика
  * repositories/ — работа с БД
  * models/ — структуры
  * middleware/ — логгер
* migrations/ — SQL-миграции
* pkg/db/ — инициализация БД
* load_tests/ — скрипты и результаты нагрузочного тестирования
* utils/ — ошибки и логгер
```
___

## ✅ Примеры запросов

**POST /api/v1/wallet**
```json
{
  "walletId": "c3a8cb84-03f2-4fb9-982a-9ee2cfb50b9f",
  "operationType": "DEPOSIT",
  "amount": 1000
}
```

**Пример ответа:**
```json
{
  "uuid": "c3a8cb84-03f2-4fb9-982a-9ee2cfb50b9f",
  "balance": 2000
}
```
___