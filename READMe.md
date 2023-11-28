# deez-nats

## Описание

Сервис, отображающий данные о заказе, полученные из канала nats-streaming

## Структура

```text
│   docker-compose.yml
│   go.mod
│   go.sum
│   Makefile
│   READMe.md
├───cmd
│   ├───publisher
│   │       main.go
│   │
│   └───subscriber
│           main.go
│
├───config
│       config.go
│
├───internal
│   ├───migrations
│   │       20231128004226_create_tables.go
│   │       migrations.go
│   │
│   ├───models
│   │       model.json
│   │       models.go
│   │
│   ├───ports
│   │   ├───publisher
│   │   │       handler.go
│   │   │       router.go
│   │   │
│   │   └───subscriber
│   │           handler.go
│   │           router.go
│   │
│   ├───repo
│   │       cached.go
│   │       permanent.go
│   │       queries.go
│   │       repo.go
│   │
│   └───service
│       ├───publisher
│       │       config.go
│       │       publisher.go
│       │
│       └───subscriber
│               config.go
│               subscriber.go
│
└───pkg
    ├───logging
    │       logging.go
    │
    └───storage
        ├───cache
        │       storage.go
        │
        └───postgres
                config.go
                storage.go
```

## Описание работы сервиса

Данные о заказе отправляются при помощи отдельно поднятого сервера для публикации данных в канал.
Для отправки данных в канал необходимо отправить POST-запрос по ссылке `http://localhost:8081/publish` 

Данные читаются и заносятся в БД так же в отдельном сервере. Данные кешируются и в случае падения сервиса
восстанавливаются из БД. Для того чтобы получить данные о заказе, необходимо отправить GET-запрос по ссылке
`http://localhost:8080/order/:id`

## Запуск

Перед запуском необходимо установить переменные окружения в файле `.env` в корне проекта.
Пример файла:

```text
HOST=localhost
PORT=8080

# db
DB_USER=user
DB_PASSWORD=password
DB_NAME=deeznats
DB_HOST=postgres
DB_PORT=5432
DB_POOL_SIZE=100

#nats-streaming
CLUSTER_ID=test-cluster
PUBLISHER_CLIENT_ID=producer
SUBSCRIBER_CLIENT_ID=subscriber
SUBJECT=subject
NATS_URL=nats://0.0.0.0:4222
```

Для запуска контейнеров с БД и сервером nats-streaming:

```bash
make up
```

Для запуска publisher'a:

```bash
make publisher
```

Для запуска subscriber'a:

```bash
make subscriber
```