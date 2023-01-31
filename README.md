# Тестовое задание 0 (стажировка Wildberries)

#### Используемые технологии:

- [Gin](https://github.com/gin-gonic/gin) - Веб фреймворк
- [SQLX](https://github.com/jmoiron/sqlx) - Библиотека для работы с базой данных
- [GoDotEnv](https://github.com/joho/godotenv) - Библиотека для работы с `.env` файлами
- [PQ](https://github.com/lib/pq) - Драйвер для PostgreSQL
- [Logrus](https://github.com/sirupsen/logrus) - Логгер
- [Viper](https://github.com/spf13/viper) - Библиотека для работы с конфигурационными файлами
- [Stan.go](https://github.com/nats-io/stan.go) - NATS Streaming System
- [Swag](https://github.com/swaggo/swag) - Автоматическое создание документации RESTful API с помощью Swagger 2.0 для Go
- [GoMock](https://github.com/golang/mock) - Mocking framework for the Go programming language.

#### Запуск:
Все сервисы можно запустить с помощью `Makefile`

`make docker` - запуск сервисов postgres, pgadmin и nats

`make swagger` - обновление данных об API для Swagger

`make publish` - добавление нового сообщения в Nats Streaming Service

`make run` - запуск основного сервера