# go-musthave-shortener-tpl
![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go)
![gRPC](https://img.shields.io/badge/-gRPC-000?logo=grpc)
![REST](https://img.shields.io/badge/-REST-FF6C37?logo=rest&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-85EA2D?logo=swagger&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-4169E1?logo=postgresql)

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы
1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона
Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:
```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:
```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.
## Запуск автотестов
Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

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
|------------|------|--------------|-----------|
| `SERVER_ADDRESS` | `-a` | `localhost:8080` | Адрес HTTP сервера |
| `SERVER_ADDRESS_GRPC` | - | `localhost:8081` | Адрес gRPC сервера |
| `URL` | `-b` | `http://localhost:8080/` | Базовый URL сервиса |
| `LOG_LEVEL` | `-l` | `INFO` | Уровень логирования (DEBUG, INFO, WARN, ERROR) |
| `FILE_STORAGE_PATH` | `-f` | `internal/storage/urls.json` | Путь к файловому хранилищу URL |
| `DATABASE_DSN` | `-d` | - | DSN для подключения к PostgreSQL |
| `AUTH_KEY` | - | hex-ключ | Ключ для аутентификации JWT |
| `TRUSTED_SUBNET` | `-t` | `192.168.1.0/24` | Доверенная подсеть (CIDR) |
| `TERMINATION_TIMEOUT` | - | `30s` | Таймаут graceful shutdown |
| `WORKER_COUNT` | - | `10` | Количество воркеров |
| `ENABLE_HTTPS` | `-s` | `false` | Включение TLS |
| `CONFIG` | `-c` | `internal/config/config.json` | Путь к JSON файлу конфигурации |