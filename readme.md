
# OrderService
Сервис для выдачи данных о заказе по его ID.

## Структура проекта

Проект организован в соответствии с best-practices для Go-приложений, с четким разделением обязанностей и удобством расширения.

Структура репозитория:

1. **cmd/**
   - точка входа приложения (`orders-service`).
   - скрипт продюсера для отправки сообщений в Kafka (`order-producer`).
  
2. **internal/**
   1. **db/**
      - Подключение к базе данных PostgreSQL (функция `ConnectDB()`).
   2. **consumer/**
      - Консьюмер Kafka (`StartKafkaConsumer`), читает сообщения и сохраняет заказы в БД.
   3. **service/**
      - Бизнес-логика: сохранение заказов (`SaveOrder`), кеширование, восстановление кеша (`RestoreCache`), доступ к данным (`GetOrderFromDB`).
   4. **server/**
      - HTTP-сервер, REST API и раздача статики.

3. **pkg/**
   - Модели данных (`Order`, `Delivery`, `Payment`, `Item`).

4. **schema/**
   - SQL-схема базы данных.

5. **static/**
   - HTML/JS страница для взаимодействия с API.

6. **docker-compose.yml**
   - Конфигурация Kafka и ZooKeeper.

7. **test.json**
   - Пример JSON-заказа для отправки в Kafka.

8. **run.sh**
   - Скрипт запуска проекта (экспорт ENV + запуск).

---

## Как запускать сервис

### Шаг 1: PostgreSQL

```bash
psql postgres <<'SQL'
CREATE DATABASE orders_service_db;
CREATE USER orders_service_user PASSWORD 'veryhardpassword12345';
GRANT ALL PRIVILEGES ON DATABASE orders_service_db TO orders_service_user;
SQL

psql -U orders_service_user -d orders_service_db -f schema/schema.sql
````

### Шаг 2: Kafka и ZooKeeper

```bash
docker compose up -d
```

### Шаг 3: Запуск сервиса

```bash
chmod +x run.sh   # один раз
./run.sh          # запуск приложения с нужными ENV
```

Ожидаемый вывод:

```
Connected PG
Cache loaded: N orders
Kafka consumer started
HTTP server on :8081
```

---

## Пример использования

### HTTP-запрос:

```bash
curl http://localhost:8081/order/b563feb7b2b84b6test | jq
```

### Ожидаемый ответ:

```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```

### Видео работы:

[Смотреть демо-видео](https://drive.google.com/file/d/1-U-Ti53Mk0OmKOQgpkMY8NvHtHTkE16J/view?usp=sharing)

---

## Продюсер 

```bash
go run ./cmd/order-producer -f test.json
```


## Технические детали

* Использован **pgxpool** для эффективной работы с PostgreSQL
* Кеш заказов реализован с помощью встроенного типа Go (`map[string]Order`)

### Ключевые решения и библиотеки

| Решение                  | Почему выбрано                                              |
| ------------------------ | ----------------------------------------------------------- |
| **segmentio/kafka-go**   | Чистый Go, удобный API, минимальные зависимости             |
| **pgxpool (jackc/pgx)**  | Высокая производительность, контекстная работа с PostgreSQL |
| **In-memory map**        | Простота, скорость чтения для демонстрации                  |
| **Транзакции в БД**      | Гарантия целостности данных                                 |
| **CommitMessages Kafka** | Подтверждение после сохранения                              |

### Архитектурные паттерны

* **Чистая архитектура** (чёткое разделение ответственности: `db`, `consumer`, `service`, `server`)
* **Идемпотентность** (`ON CONFLICT` SQL-запросы, безопасны при повторе)
* **Явное подтверждение обработки сообщений Kafka**

### Основные компоненты и пакеты:

* Подключение к БД: [`internal/db`](internal/db)
* Бизнес-логика и кеш: [`internal/service`](internal/service)
* Консьюмер Kafka: [`internal/consumer`](internal/consumer)
* HTTP API: [`internal/server`](internal/server)
* Структуры данных: [`pkg/models`](pkg/models)

---

