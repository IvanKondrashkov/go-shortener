# go-musthave-shortener-tpl
![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go)
![gRPC](https://img.shields.io/badge/-gRPC-000?logo=grpc)
![REST](https://img.shields.io/badge/-REST-FF6C37?logo=rest&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-85EA2D?logo=swagger&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-4169E1?logo=postgresql)

## Как запустить контейнер
Сборка бинарных файлов сервера:
```shell
task build-server
```

Запустите локально Docker:
```shell
docker-compose up -d
```

## Конфигурация
### Конфигурация Сервера
Сервер поддерживает настройку через переменные окружения, аргументы командной строки и JSON файл:

| Переменная | Флаг | По умолчанию | Описание |
|------------|---|--------------|-----------|
| `SERVER_ADDRESS` | `-a` | `localhost:8080` | Адрес HTTP сервера |
| `SERVER_ADDRESS_GRPC` | - | `localhost:8081` | Адрес gRPC сервера |
| `URL` | `-b` | `http://localhost:8080/` | Базовый URL сервиса |
| `LOG_LEVEL` | `-l` | `INFO` | Уровень логирования (DEBUG, INFO, WARN, ERROR) |
| `FILE_STORAGE_PATH` | `-f` | `internal/storage/urls.json` | Путь к файловому хранилищу URL |
| `DATABASE_DSN` | `-d` | - | DSN для подключения к PostgreSQL |
| `AUTH_KEY` | - | hex-ключ | Ключ для аутентификации JWT |
| `TRUSTED_SUBNET` | `-t` | `192.168.1.0/24` | Доверенная подсеть (CIDR) |
| `PATH_KEY` | - | `cert/server.key` | Путь до ключа TLS |
| `PATH_CERT` | - | `cert/server.crt` | Путь до сертификата TLS |
| `TERMINATION_TIMEOUT` | - | `30s` | Таймаут graceful shutdown |
| `WORKER_COUNT` | - | `10` | Количество воркеров |
| `ENABLE_HTTPS` | `-s` | `false` | Включение TLS |
| `CONFIG` | `-c` | `internal/config/config.json` | Путь к JSON файлу конфигурации |