# Тестовое задание 0 (стажировка Wildberries)

#### Используемые технологии:

- [Gin](https://github.com/gin-gonic/gin) - Веб фреймворк
- [SQLX](https://github.com/jmoiron/sqlx) - Библиотека для работы с базой данных
- [GoDotEnv](https://github.com/joho/godotenv) - Библиотека для работы с `.env` файлами
- [PQ](github.com/lib/pq) - Драйвер для PostgreSQL
- [Logrus](github.com/sirupsen/logrus) - Логгер
- [Viper](https://github.com/spf13/viper) - Библиотека для работы с конфигурационными файлами

#### Запуск:
Сервисы `postgres` и `pgadmin` запускаются с помощью `docker-compose`
Сервер запускается отдельно, например, с помощью `Makefile`